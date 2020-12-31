package res

import (
	"encoding/json"
	"net/http"
)

// JSON - send a json response
func JSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error", 500)
	}
}

// Plain - send a plain text
func Plain(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "text/plain")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error", 500)
	}
}
