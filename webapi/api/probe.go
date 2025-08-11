package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ttrnecka/agent_poc/webapi/db"
	"github.com/ttrnecka/agent_poc/webapi/ws"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) ProbeApiHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		probes, err := db.Probes().CRUD().All(r.Context())
		// probes, err := db.GetProbes(r.Context())
		if err != nil {
			logger.Error().Err(err).Msg("")
			http.Error(w, "Failed to fetch probes", http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(probes); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		p := r.URL.Path
		var probe_id string
		switch {
		case match(p, "/api/v1/probe/+", &probe_id), p == "/api/v1/probe":
			var probe db.Probe
			if err := json.NewDecoder(r.Body).Decode(&probe); err != nil {
				http.Error(w, fmt.Sprintf("Invalid JSON: %s", err), http.StatusBadRequest)
				return
			}
			id, err := probe.UpdateProbe(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
				return
			}
			probeTmp, err := db.Probes().CRUD().GetByID(r.Context(), id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
				return
			}
			json.NewEncoder(w).Encode(probeTmp)
			go h.refreshPolicies()
			return
		default:
			http.Error(w, "Unhandled path", http.StatusInternalServerError)
			return
		}
	case http.MethodDelete:
		p := r.URL.Path
		var probe_id string
		switch {
		case match(p, "/api/v1/probe/+", &probe_id):
			id, err := primitive.ObjectIDFromHex(probe_id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
				return
			}
			probe, err := db.Probes().CRUD().GetByID(r.Context(), id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
				return
			}
			err = probe.Delete(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
				return
			}
		default:
			fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func (h *Handler) refreshPolicies() {
	probes, err := db.Probes().CRUD().All(context.Background())
	if err != nil {
		logger.Error().Err(err).Msg("Refreshing policies")
		return
	}

	collectors := make(map[string]bool)

	for _, p := range probes {
		collectors[p.Collector().Name] = true
	}

	// once saved we broadcast the new probes to all connected clients
	hub := ws.GetHub()
	for collector := range collectors {
		bmessage, err := json.Marshal(ws.NewMessage(ws.MSG_POLICY_REFRESH, "hub", collector, "Policy updated"))
		if err != nil {
			logger.Error().Err(err).Msg("Refreshing policies, marshall message")
			continue
		}
		hub.BroadcastMessage(bmessage)
	}
}
