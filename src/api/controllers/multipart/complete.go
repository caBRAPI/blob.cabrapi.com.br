package multipart

import (
	"blob/src/auth"
	"blob/src/services"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// CompleteUpload handles POST /blob/{uploadId}/complete
func CompleteUpload(w http.ResponseWriter, r *http.Request) {
	finalHash := r.Header.Get("X-Final-Hash")
	if finalHash == "" {
		writeJSONError(w, "Missing X-Final-Hash header", http.StatusBadRequest)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		writeJSONError(w, "Missing uploadId", http.StatusBadRequest)
		return
	}
	uploadID, err := uuid.Parse(parts[2])
	if err != nil {
		writeJSONError(w, "Invalid uploadId", http.StatusBadRequest)
		return
	}
	userID := auth.UserIDFromRequest(r)
	if userID == uuid.Nil {
		writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	blob, err := services.Multipart.Complete(r.Context(), services.CompleteInput{
		UploadID:  uploadID,
		UserID:    userID,
		FinalHash: finalHash,
	})
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, services.ErrUnauthorized) {
			status = http.StatusForbidden
		} else if strings.Contains(err.Error(), "invalid uploadId") {
			status = http.StatusNotFound
		}
		writeJSONError(w, err.Error(), status)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":         blob.ID.String(),
		"bucket":     blob.Bucket,
		"filename":   blob.Filename,
		"size":       blob.Size,
		"hash":       blob.Hash,
		"created_at": blob.CreatedAt,
		"updated_at": blob.UpdatedAt,
	})
}
