package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"blob/src/database"
	"blob/src/functions"
	"blob/src/models"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func writeJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func UploadBlobController(w http.ResponseWriter, r *http.Request) {
	// Parse expires_at if provided
	expiresAtStr := r.FormValue("expires_at")
	var expiresAt *time.Time
	if expiresAtStr != "" {
		t, err := time.Parse(time.RFC3339, expiresAtStr)
		if err == nil {
			expiresAt = &t
		}
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeJSONError(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSONError(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	bucket := r.FormValue("bucket")
	validate := validator.New()
	input := struct {
		Bucket string `validate:"required,min=1,max=64"`
	}{Bucket: bucket}
	if err := validate.Struct(input); err != nil {
		writeJSONError(w, "Bucket is required", http.StatusBadRequest)
		return
	}

	allowedMimes := os.Getenv("BLOB_ALLOWED_MIME_TYPES")
	maxUploadSize, _ := strconv.ParseInt(os.Getenv("BLOB_MAX_UPLOAD_SIZE_BYTES"), 10, 64)
	maxStorageSize, _ := strconv.ParseInt(os.Getenv("BLOB_MAX_STORAGE_SIZE"), 10, 64)
	storagePath := os.Getenv("BLOB_STORAGE_PATH")
	if storagePath == "" {
		storagePath = "storage/uploads"
	}

	mime := header.Header.Get("Content-Type")
	if !functions.IsAllowedMimeType(mime, functions.SplitComma(allowedMimes)) {
		writeJSONError(w, "MIME type not allowed", http.StatusBadRequest)
		return
	}

	if maxUploadSize > 0 && header.Size > maxUploadSize {
		writeJSONError(w, "File too large", http.StatusRequestEntityTooLarge)
		return
	}

	if maxStorageSize > 0 {
		total, err := functions.GetTotalStorageSize(storagePath)
		if err == nil && total+header.Size > maxStorageSize {
			writeJSONError(w, "Storage limit exceeded", http.StatusInsufficientStorage)
			return
		}
	}

	id := uuid.New()
	// Permitir filename customizado
	filename := r.FormValue("filename")
	if filename == "" {
		filename = header.Filename
	}
	blobPath := storagePath + string(os.PathSeparator) + id.String()
	os.MkdirAll(storagePath, 0755)
	out, err := os.Create(blobPath)
	if err != nil {
		writeJSONError(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	size, err := io.Copy(out, file)
	if err != nil {
		writeJSONError(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	metadata := r.FormValue("metadata")
	var metaJson map[string]interface{}
	if metadata != "" {
		json.Unmarshal([]byte(metadata), &metaJson)
	}

	hash := id.String()

	// Permitir campo public customizado
	public := true
	publicStr := r.FormValue("public")
	if publicStr != "" {
		if publicStr == "false" || publicStr == "0" {
			public = false
		}
	}
	blob := models.Blob{
		ID:        id,
		Bucket:    bucket,
		Filename:  filename,
		Mime:      mime,
		Size:      size,
		Hash:      hash,
		Path:      blobPath,
		Public:    &public,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}
	if metaJson != nil {
		metaBytes, _ := json.Marshal(metaJson)
		blob.Metadata = metaBytes
	}

	if err := database.DB.Create(&blob).Error; err != nil {
		writeJSONError(w, "Failed to save blob metadata", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Monta a URL de download usando o host do request
	url := ""
	if r.TLS != nil {
		url = "https://" + r.Host + "/blob/" + blob.ID.String()
	} else {
		url = "http://" + r.Host + "/blob/" + blob.ID.String()
	}

	type responseWithURL struct {
		models.Blob
		URL string `json:"url"`
	}

	resp := responseWithURL{
		Blob: blob,
		URL:  url,
	}

	json.NewEncoder(w).Encode(resp)
}
