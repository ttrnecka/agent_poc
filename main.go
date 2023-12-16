package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	ui "github.com/ttrnecka/agent_poc/agent_poc"
	"github.com/ttrnecka/agent_poc/api"
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
		next.ServeHTTP(w, r)
	})
}

func router() http.Handler {
	mux := http.NewServeMux()

	// index page
	mux.HandleFunc("/", indexHandler)

	// static files
	staticFS, _ := fs.Sub(ui.StaticFiles, "dist")
	httpFS := http.FileServer(http.FS(staticFS))
	mux.Handle("/assets/", httpFS)

	// api
	mux.HandleFunc("/api/v1/greeting", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, there!"))
	})

	mux.HandleFunc("/api/v1/policy", commonApiMiddleware(api.PolicyApiHandler))
	mux.HandleFunc("/api/v1/probe", commonApiMiddleware(api.ProbeApiHandler))

	return mux
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api") {
		http.NotFound(w, r)
		return
	}

	if r.URL.Path == "/favicon.ico" {
		rawFile, _ := ui.StaticFiles.ReadFile("dist/favicon.ico")
		w.Write(rawFile)
		return
	}
	rawFile, _ := ui.StaticFiles.ReadFile("dist/index.html")
	w.Write(rawFile)
}
