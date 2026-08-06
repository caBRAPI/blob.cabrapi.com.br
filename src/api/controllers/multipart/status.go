package multipart

import (
	"blob/src/services"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// UploadStatus handles GET /blob/{uploadId}/status
func UploadStatus(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Missing uploadId", http.StatusBadRequest)
		return
	}
	uploadID, err := uuid.Parse(parts[2])
	if err != nil {
		http.Error(w, "Invalid uploadId", http.StatusBadRequest)
		return
	}
	userID, err := uuid.Parse(r.Header.Get("X-User-ID"))
	if err != nil {
		http.Error(w, "Missing or invalid X-User-ID header", http.StatusUnauthorized)
		return
	}

	chunkSize := 0
	if cs := r.URL.Query().Get("chunk_size"); cs != "" {
		if v, err := strconv.Atoi(cs); err == nil && v > 0 {
			chunkSize = v
		}
	}

	status, err := services.Multipart.Status(r.Context(), uploadID, userID, chunkSize)
	if err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, services.ErrUnauthorized) {
			statusCode = http.StatusForbidden
		} else if strings.Contains(err.Error(), "invalid uploadId") {
			statusCode = http.StatusNotFound
		}
		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
