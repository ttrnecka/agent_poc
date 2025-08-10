package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ttrnecka/agent_poc/webapi/db"
	"go.mongodb.org/mongo-driver/bson"
)

type Policy struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
}

func (h *Handler) PolicyApiHandler(w http.ResponseWriter, r *http.Request) {
	policies, err := db.Policies().Find(r.Context(), bson.D{})
	if err != nil {
		http.Error(w, "Failed to fetch policies", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(policies); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

func PolicyItemApiHandler(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	var policy_name, policy_version string
	switch {
	case match(p, "/api/v1/policy/+/+", &policy_name, &policy_version):
		file := fmt.Sprintf("policies/bin/%s_%s", policy_name, policy_version)
		handleLargeFile(w, r, file)
	default:
		http.NotFound(w, r)
		return
	}
}

func handleLargeFile(w http.ResponseWriter, r *http.Request, LargeFileName string) {
	//Open file
	f, err := os.Open(LargeFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, err.Error(), 404)
			return
		}
		http.Error(w, err.Error(), 400)
		return
	}
	defer f.Close()

	//read the file info
	info, err := f.Stat()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	//Set the headers
	setHeaders(w, filepath.Base(LargeFileName), strconv.Itoa(int(info.Size())))
	w.WriteHeader(http.StatusOK)

	//Copy without loading everything in memory
	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func setHeaders(w http.ResponseWriter, name, len string) {
	//Represents binary file
	w.Header().Set("Content-Type", "application/octet-stream")
	//Tells client what filename should be used.
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, name))
	//The length of the data.
	w.Header().Set("Content-Length", len)
	//No cache headers.
	w.Header().Set("Cache-Control", "private")
	//No cache headers.
	w.Header().Set("Pragma", "private")
	//No cache headers.
	w.Header().Set("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
}
