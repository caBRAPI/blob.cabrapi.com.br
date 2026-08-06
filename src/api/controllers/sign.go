package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"blob/src/auth"
	"blob/src/functions"
	"blob/src/services"
)

// SignBlobController handles POST /blob/{id}/sign (private).
// It returns expiring download and view URLs signed with HMAC-SHA256.
func SignBlobController(w http.ResponseWriter, r *http.Request) {
	id, ok := parseBlobID(r)
	if !ok {
		functions.WriteJSONError(w, "Invalid blob id", http.StatusBadRequest)
		return
	}

	if _, err := services.Blobs.Get(r.Context(), id); err != nil {
		functions.WriteJSONError(w, "Blob not found", http.StatusNotFound)
		return
	}

	expires := time.Now().UTC().Add(auth.DefaultTTL())
	expiresUnix := expires.Unix()

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	base := scheme + "://" + r.Host

	downloadSig := auth.Sign(auth.Secret(), id.String(), auth.ActionDownload, expires)
	viewSig := auth.Sign(auth.Secret(), id.String(), auth.ActionView, expires)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         id.String(),
		"expires_at": expires.Format(time.RFC3339),
		"ttl":        int64(auth.DefaultTTL().Seconds()),
		"download":   base + "/blob/" + id.String() + "/download?expires=" + strconv.FormatInt(expiresUnix, 10) + "&signature=" + downloadSig,
		"view":       base + "/blob/" + id.String() + "/view?expires=" + strconv.FormatInt(expiresUnix, 10) + "&signature=" + viewSig,
	}); err != nil {
		functions.Error("failed to encode sign json: %v", err)
	}
}
