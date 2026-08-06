package routes

import (
	"encoding/json"
	"net/http"

	"blob/src/version"
)

// VersionHandler handles GET /version
func VersionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"version": version.V,
	})
}
