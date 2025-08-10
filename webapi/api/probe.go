package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ttrnecka/agent_poc/webapi/db"
	"github.com/ttrnecka/agent_poc/webapi/ws"
)

func (h *Handler) ProbeApiHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		probes, err := db.GetProbes(r.Context())
		if err != nil {
			logger.Error().Err(err).Msg("")
			http.Error(w, "Failed to fetch probes", http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(probes); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}
	case "POST":
		if err := h.saveProbes(r.Body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Save err : %v", err)
			return
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func (h *Handler) saveProbes(r io.Reader) error {
	// outputFileName := "data/api/probes.json"

	var result []interface{}
	err := json.NewDecoder(r).Decode(&result)
	if err != nil {
		return err
	}

	err = db.SaveProbes(result)
	if err != nil {
		return err
	}

	probes, err := db.GetProbes(context.Background())
	if err != nil {
		return err

	}

	collectors := make(map[string]bool)

	for _, p := range probes {
		collectors[p.Collector] = true
	}

	// once saved we broadcast the new probes to all connected clients
	hub := ws.GetHub()
	for collector := range collectors {
		bmessage, err := json.Marshal(ws.NewMessage(ws.MSG_POLICY_REFRESH, "hub", collector, "Policy updated"))
		if err != nil {
			return fmt.Errorf("failed to marshal message: %v", err)
		}
		hub.BroadcastMessage(bmessage)
	}

	return nil
}
