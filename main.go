package main

import (
	"blob/src/database"
	"blob/src/functions"
	"blob/src/middleware"
	"blob/src/routes"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	database.Redis()
	database.Postgres()

	mux := http.NewServeMux()
	limiter := middleware.NewRateLimiterFromEnv()
	routes.RegisterRoutes(mux, limiter)

	PORT := os.Getenv("BLOB_PORT")
	if PORT == "" {
		PORT = "3000"
	}
	HOST := os.Getenv("BLOB_HOST")
	if HOST == "" {
		HOST = "localhost"
	}

	now := time.Now().Format("Mon Jan 2 15:04:05 2006")

	functions.Info("[SERVER] Server running at: http://%s:%s (%s)", HOST, PORT, now)
	http.ListenAndServe(HOST+":"+PORT, mux)
}
