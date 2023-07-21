package middleware

import (
	"log"
	"net/http"
	"time"
)

// The LoggingMiddleware function is a middleware that logs the HTTP method, URL, and the time it took
// to process the request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(w, r)

		log.Printf("%s %s %d us", r.Method, r.URL, time.Since(startTime).Microseconds())
	})
}
