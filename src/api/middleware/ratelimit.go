package middleware

import (
	"blob/src/config"
	"blob/src/database"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"
)

type RateLimiter struct {
	max    int
	window time.Duration
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func Variables() *RateLimiter {
	max := config.Env.RateLimitMax
	window := time.Duration(config.Env.RateLimitWindowMS) * time.Millisecond
	return &RateLimiter{
		max:    max,
		window: window,
	}
}

// clientIP resolves the real client IP, honoring forwarded headers only when
// the proxy is trusted (BLOB_TRUST_PROXY=true). Otherwise RemoteAddr is used.
func clientIP(r *http.Request) string {
	if config.Env != nil && config.Env.TrustProxy {
		if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
			first := strings.TrimSpace(strings.Split(fwd, ",")[0])
			if first != "" {
				return first
			}
		}
		if real := r.Header.Get("X-Real-IP"); real != "" {
			return strings.TrimSpace(real)
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		key := "ratelimit:" + ip
		count, err := database.RedisClient.Incr(database.Ctx, key).Result()
		if err == nil && count == 1 {
			database.RedisClient.PExpire(database.Ctx, key, rl.window)
		}
		if err != nil || int(count) > rl.max {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "rate limit exceeded"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
