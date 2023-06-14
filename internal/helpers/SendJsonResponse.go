package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SendJsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	dataJson, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%s", dataJson)
}
