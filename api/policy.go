package api

import (
	"fmt"
	"net/http"
)

func PolicyApiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, output("policies.json"))
}
