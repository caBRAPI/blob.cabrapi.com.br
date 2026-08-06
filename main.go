package main

import (
	"blob/src/api/middleware"
	"blob/src/api/routes"
	"blob/src/config"
	"blob/src/database"
	"blob/src/functions"
	"blob/src/services"
	queue "blob/src/services/queue"
	"blob/src/version"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var Version = "Desenvolvimento"

func main() {

	functions.Info("Starting Blob API Version: %s", Version)

	if err := godotenv.Load(); err != nil {
		functions.Error("No .env file found, relying on environment variables")
	}

	database.Redis()
	database.Postgres()

	cfg := config.Load()
	version.V = Version

	if cfg.SignedURLSecret == "" || len(cfg.SignedURLSecret) < 16 {
		functions.Error("BLOB_SIGNED_URL_SECRET is required (at least 16 characters). Signed URLs are disabled until it is set.")
		os.Exit(1)
	}

	services.InitAsynq()
	services.InitBlobService()
	services.InitMultipartService()
	queue.StartQueueWorker()
	queue.StartCleanupScheduler()
	queue.StartTmpCleanupScheduler()

	mux := http.NewServeMux()
	limiter := middleware.Variables()
	routes.RegisterRoutes(mux, limiter)

	corsOpts := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Accept", "Origin", "X-User-ID", "X-Chunk-Index", "X-Chunk-Hash", "X-Final-Hash", "Range"},
	}
	if cfg.CORSOrigins != "*" && cfg.CORSOrigins != "" {
		origins := strings.Split(cfg.CORSOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}
		corsOpts.AllowedOrigins = origins
	}
	handler := cors.New(corsOpts).Handler(mux)

	functions.Info("[SERVER] Server running at: http://%s:%s", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		functions.Error("server error: %v", err)
	}

}
