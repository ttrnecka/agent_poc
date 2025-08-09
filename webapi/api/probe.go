package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/ttrnecka/agent_poc/webapi/ws"
)

type Probe struct {
	Id        string `json:"id"`
	Collector string `json:"collector"`
	Policy    string `json:"policy"`
	Version   string `json:"version"`
	Address   string `json:"address"`
	Port      string `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
}

func (h *Handler) ProbeApiHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		probes, err := h.DB.GetAllProbes(r.Context())
		if err != nil {
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

	// // Create or open the output file for writing.
	// outputFile, err := os.Create(outputFileName)
	// if err != nil {
	// 	return err
	// }
	// defer outputFile.Close()

	var probes []Probe
	err := json.NewDecoder(r).Decode(&probes)
	if err != nil {
		return err
	}

	collectors := make(map[string]bool)

	for i, p := range probes {
		collectors[p.Collector] = true
		if p.Id == "" {
			probes[i].Id = uuid.New().String()
		}
	}

	result := make([]interface{}, len(probes))
	for i, p := range probes {
		result[i] = p
	}

	err = h.DB.SaveProbes(result)
	if err != nil {
		return err
	}

	// var buf bytes.Buffer
	// err = json.NewEncoder(&buf).Encode(probes)
	// if err != nil {
	// 	return err
	// }

	// // Copy the contents of the Reader to the output file.
	// _, err = io.Copy(outputFile, &buf)
	// if err != nil {
	// 	return err
	// }

	// once saved we broadcast the new probes to all connected clients
	hub := ws.GetHub()
	for collector := range collectors {
		bmessage, err := json.Marshal(ws.NewMessage(ws.MSG_REFRESH, "hub", collector, "Policy updated"))
		if err != nil {
			return fmt.Errorf("failed to marshal message: %v", err)
		}
		hub.BroadcastMessage(bmessage)
	}

	return nil
}
