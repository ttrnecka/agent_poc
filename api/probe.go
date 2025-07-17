package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
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

func ProbeApiHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, output("probes.json"))
	case "POST":

		if err := save(r.Body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Save err : %v", err)
			return
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

// func loadProbe(uuid string) (*Probe, error) {
// 	data := output("probes.json")
// 	var unm []Probe
// 	err := json.Unmarshal([]byte(data), &unm)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, p := range unm {
// 		if p.Id == uuid {
// 			return &p, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("Probe with uuid %s not found", uuid)
// }

func save(r io.Reader) error {
	outputFileName := "data/probes.json"

	// Create or open the output file for writing.
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	var unm []Probe
	err = json.NewDecoder(r).Decode(&unm)
	if err != nil {
		return err
	}
	for i, p := range unm {
		if p.Id == "" {
			unm[i].Id = uuid.New().String()
		}
	}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(unm)
	if err != nil {
		return err
	}

	// Copy the contents of the Reader to the output file.
	_, err = io.Copy(outputFile, &buf)
	if err != nil {
		return err
	}

	return nil
}
