package api

import (
	"fmt"
	"net/http"
)

func CollectorApiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, output("collectors.json"))
}
