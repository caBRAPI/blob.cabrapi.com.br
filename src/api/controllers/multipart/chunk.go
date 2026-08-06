package multipart

import (
	"blob/src/auth"
	"blob/src/config"
	"blob/src/services"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// UploadChunk handles PUT /blob/{uploadId}/chunk
func UploadChunk(w http.ResponseWriter, r *http.Request) {
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
	userID := auth.UserIDFromRequest(r)
	if userID == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chunkIdxStr := r.Header.Get("X-Chunk-Index")
	if chunkIdxStr == "" {
		http.Error(w, "Missing X-Chunk-Index header", http.StatusBadRequest)
		return
	}
	var chunkIdx int
	if _, err := fmt.Sscanf(chunkIdxStr, "%d", &chunkIdx); err != nil || chunkIdx < 0 {
		http.Error(w, "Invalid X-Chunk-Index header", http.StatusBadRequest)
		return
	}

	cfg := config.Env
	r.Body = http.MaxBytesReader(w, r.Body, cfg.MaxChunkSize)

	err = services.Multipart.UploadChunk(r.Context(), services.ChunkInput{
		UploadID: uploadID,
		UserID:   userID,
		Index:    chunkIdx,
		Hash:     r.Header.Get("X-Chunk-Hash"),
		MaxSize:  cfg.MaxChunkSize,
		Body:     r.Body,
	})
	if err != nil {
		writeChunkError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeChunkError(w http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, services.ErrUnauthorized) {
		status = http.StatusForbidden
	} else if strings.Contains(err.Error(), "invalid uploadId") {
		status = http.StatusNotFound
	}
	http.Error(w, err.Error(), status)
}
