package api

import (
	"log"
	"net/http"
)

// LogClientIP logs the client IP for each request before passing it on.
func LogClientIP(next http.Handler) http.Handler {
	return logClientIP(next, log.Default())
}

func logClientIP(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ip := clientIP(r); ip != "" {
			logger.Printf("client_ip=%s", ip)
		}
		next.ServeHTTP(w, r)
	})
}
