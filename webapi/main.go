package main

import (
	"fmt"
	"net/http"

	"github.com/ttrnecka/agent_poc/webapi/api"
	"github.com/ttrnecka/agent_poc/webapi/ws"
)

func main() {
	srv := &http.Server{
		Addr:    ":8888",
		Handler: router(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func commonApiMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		next.ServeHTTP(w, r)
	})
}

func router() http.Handler {
	mux := http.NewServeMux()

	// index page
	mux.HandleFunc("/", indexHandler)

	// api
	mux.HandleFunc("/api/v1/policy", commonApiMiddleware(api.PolicyApiHandler))
	mux.HandleFunc("/api/v1/policy/", commonApiMiddleware(api.PolicyItemApiHandler))
	mux.HandleFunc("/api/v1/probe", commonApiMiddleware(api.ProbeApiHandler))
	mux.HandleFunc("/api/v1/collector", commonApiMiddleware(api.CollectorApiHandler))
	mux.HandleFunc("/api/v1/data/collector/", commonApiMiddleware(api.DataHandler))
	mux.HandleFunc("/api/v1/data/collector", commonApiMiddleware(api.DataHandler)) // handles no trailing slash too

	hub := ws.GetHub()
	go hub.Run()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	return mux
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
