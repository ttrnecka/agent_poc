package api

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) CollectorApiHandler(w http.ResponseWriter, r *http.Request) {

	collectors, err := h.DB.GetAllCollectors(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch collectors", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(collectors); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}
