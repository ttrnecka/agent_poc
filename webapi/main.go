package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/rs/zerolog"
	"github.com/ttrnecka/agent_poc/webapi/api"
	"github.com/ttrnecka/agent_poc/webapi/ws"
	"golang.org/x/crypto/bcrypt"

	logging "github.com/ttrnecka/agent_poc/logger"
)

var logger zerolog.Logger

func init() {
	logger = logging.SetupLogger("webapi")
}

var (
	sessionManager *scs.SessionManager
	users          = map[string]User{} // simple in-memory user store
)

func main() {

	// Needed for storing structs in sessions
	gob.Register(User{})

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Secure = false // Set to true in production

	// Preload a test user
	hashedPw, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	users["test"] = User{ID: 1, Email: "test@test.com", Username: "test", Password: string(hashedPw)}

	srv := &http.Server{
		Addr:    ":8888",
		Handler: router(),
	}

	err := srv.ListenAndServe()
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

func router() http.Handler {
	// mux := http.NewServeMux()
	r := NewMiddlewareRouter()
	//api
	r.Handle("/api/v1/policy", http.HandlerFunc(api.PolicyApiHandler))
	r.Handle("/api/v1/policy/", http.HandlerFunc(api.PolicyItemApiHandler))
	r.Handle("/api/v1/probe", http.HandlerFunc(api.ProbeApiHandler))
	r.Handle("/api/v1/collector", http.HandlerFunc(api.CollectorApiHandler))
	r.Handle("/api/v1/data/collector/", http.HandlerFunc(api.DataHandler))
	r.Handle("/api/v1/data/collector", http.HandlerFunc(api.DataHandler)) // handles no trailing slash too

	// index
	r.Handle("/", http.HandlerFunc(indexHandler))
	r.Handle("/api/login", http.HandlerFunc(loginHandler))
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

// no index
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	http.NotFound(w, r)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Println("username", username, "passowd", password)
	user, ok := users[username]
	if !ok || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
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

	user, ok := sessionManager.Get(r.Context(), "user").(User)
	if !ok {
		http.Error(w, "type assertion to User failed", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(user.GetResponse())
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
