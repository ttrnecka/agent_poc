package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

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

func save(r io.Reader) error {
	outputFileName := "data/probes.json"

	// Create or open the output file for writing.
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Copy the contents of the Reader to the output file.
	_, err = io.Copy(outputFile, r)
	if err != nil {
		return err
	}

	return nil
}
