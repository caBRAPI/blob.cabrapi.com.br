package multipart

import (
	"encoding/json"
	"net/http"
)

// writeJSONError writes a JSON error response.
func writeJSONError(w http.ResponseWriter, message string, code int) {
	writeJSON(w, code, map[string]string{"error": message})
}

// writeJSON writes a JSON body with the given status.
func writeJSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}
