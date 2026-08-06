package multipart

import (
	"encoding/json"
	"net/http"

	"blob/src/api/validators"
	"blob/src/functions"
	"blob/src/services"

	"github.com/google/uuid"
)

// InitiateUpload handles POST /blob/initiate
func InitiateUpload(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Bucket   string `json:"bucket"`
		Filename string `json:"filename"`
		Size     int64  `json:"size"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Bucket == "" || req.Filename == "" || req.Size <= 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if !validators.ValidateMultipartBucket(req.Bucket) {
		http.Error(w, "Invalid bucket name", http.StatusBadRequest)
		return
	}
	userID, err := uuid.Parse(r.Header.Get("X-User-ID"))
	if err != nil {
		http.Error(w, "Missing or invalid X-User-ID header", http.StatusUnauthorized)
		return
	}

	upload, err := services.Multipart.Initiate(r.Context(), services.InitiateInput{
		UserID:   userID,
		Bucket:   req.Bucket,
		Filename: req.Filename,
		Size:     req.Size,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(map[string]string{"uploadId": upload.ID.String()}); err != nil {
		functions.Error("failed to encode uploadId json: %v", err)
	}
}
