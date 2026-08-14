package controllers

import (
	"net/http"
	"strings"

	"blob/src/auth"
	"blob/src/functions"
)

// ViewBlobController handles GET /blob/{id}/view
func ViewBlobController(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 || parts[2] != "view" {
		functions.WriteJSONError(w, "Invalid view URL", http.StatusBadRequest)
		return
	}

	serveBlobFile(w, r, auth.ActionView, "inline", "view_count")
}
