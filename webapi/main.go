package main

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/rs/zerolog"
	"github.com/ttrnecka/agent_poc/webapi/api"
	"github.com/ttrnecka/agent_poc/webapi/db"
	"github.com/ttrnecka/agent_poc/webapi/ws"
	"golang.org/x/crypto/bcrypt"

	cdb "github.com/ttrnecka/agent_poc/common/db"
	logging "github.com/ttrnecka/agent_poc/logger"
)

var logger zerolog.Logger

func init() {
	logger = logging.SetupLogger("webapi")
}

var (
	sessionManager *scs.SessionManager
)

func main() {

	// Needed for storing structs in sessions
	gob.Register(db.User{})

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Secure = false // Set to true in production

	// db

	dB, err := db.Connect()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	// Ensure all indexes before starting application logic
	if err := db.EnsureUserCollection(dB); err != nil {
		log.Fatal("Failed to ensure indexes:", err)
	}

	srv := &http.Server{
		Addr:    ":8888",
		Handler: router(dB),
	}

	err = srv.ListenAndServe()
	if err != nil {
		logger.Error().Err(err).Msg("")
	}
}

func commonApiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		next.ServeHTTP(w, r)
	})
}

func router(dB *db.DB) http.Handler {
	// mux := http.NewServeMux()
	r := NewMiddlewareRouter()

	handler := api.NewHandler(dB)
	cHandler := NewCoreHandler(dB)

	//api
	r.Handle("/api/v1/policy", http.HandlerFunc(handler.PolicyApiHandler))
	r.Handle("/api/v1/policy/", http.HandlerFunc(api.PolicyItemApiHandler))
	r.Handle("/api/v1/probe", http.HandlerFunc(handler.ProbeApiHandler))
	r.Handle("/api/v1/collector", http.HandlerFunc(handler.CollectorApiHandler))
	r.Handle("/api/v1/data/collector/", http.HandlerFunc(api.DataHandler))
	r.Handle("/api/v1/data/collector", http.HandlerFunc(api.DataHandler)) // handles no trailing slash too

	// index
	r.Handle("/", http.HandlerFunc(cHandler.indexHandler))
	r.Handle("/api/login", http.HandlerFunc(cHandler.loginHandler))
	r.Handle("/api/user", http.HandlerFunc(userHandler))
	r.Handle("/api/logout", http.HandlerFunc(logoutHandler))

	// Middleware definitions
	r.UseForPrefix("/api/v1/", authMiddleware, commonApiMiddleware, sessionManager.LoadAndSave)
	r.UseForPrefix("/api/login", commonApiMiddleware, sessionManager.LoadAndSave)
	r.UseForPrefix("/api/user", authMiddleware, commonApiMiddleware, sessionManager.LoadAndSave)
	r.UseForPrefix("/api/logout", commonApiMiddleware, sessionManager.LoadAndSave)
	r.UseDefault(sessionManager.LoadAndSave) // applies to /public, etc.

	hub := ws.GetHub()
	go hub.Run()
	r.Handle("/ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	}))

	return r
}

type coreHandler struct {
	DB *db.DB
}

func NewCoreHandler(db *db.DB) *coreHandler {
	return &coreHandler{DB: db}
}

// no index
func (c *coreHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	http.NotFound(w, r)
}

func (c *coreHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	ctx := context.Background()
	userCRUD := cdb.NewCRUD[db.User](c.DB.Database(), "users")
	user, err := userCRUD.GetByField(ctx, "username", username)

	// TODO add user not found check
	if err != nil {
		logger.Error().Err(err).Msg("")
		if errors.Is(err, cdb.ErrNotFound) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	sessionManager.Put(r.Context(), "user", user)

	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful", "user": user.Username})
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := sessionManager.Get(r.Context(), "user").(db.User)
	if !ok {
		http.Error(w, "type assertion to User failed", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(user)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionManager.Destroy(r.Context())
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sessionManager.Exists(r.Context(), "user") {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}
