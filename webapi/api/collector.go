package api

import (
	"encoding/json"
	"net/http"

	"github.com/ttrnecka/agent_poc/webapi/db"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *Handler) CollectorApiHandler(w http.ResponseWriter, r *http.Request) {

	collectors, err := db.Collectors().Find(r.Context(), bson.D{})
	if err != nil {
		http.Error(w, "Failed to fetch collectors", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(collectors); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}
