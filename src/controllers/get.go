package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"blob/src/database"
	"blob/src/models"

	"github.com/google/uuid"
)

// GetBlobController handles GET /blob/{id}
func GetBlobController(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing blob id"})
		return
	}
	idStr := parts[2]
	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid blob id"})
		return
	}

	var blob models.Blob
	result := database.DB.First(&blob, "id = ?", id)
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Blob not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blob)
}
