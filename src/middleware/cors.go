package middleware

import (
	"net/http"
	"os"
	"strings"
)

func CORSMiddleware(next http.Handler) http.Handler {
	allowed := os.Getenv("BLOB_CORS_ORIGINS")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allowed == "*" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			origins := strings.Split(allowed, ",")
			origin := r.Header.Get("Origin")
			for _, o := range origins {
				if strings.TrimSpace(o) == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS,HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,Accept,Origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
