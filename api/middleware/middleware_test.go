package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	// Wrap the test handler with the logging middleware
	handler := LoggingMiddleware(testHandler)

	// Create a buffer to capture the log output
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the log output
	logOutput := buf.String()
	expectedLogOutput := "GET /test"
	if !strings.Contains(logOutput, expectedLogOutput) {
		t.Errorf("handler logged wrong output: got %v want %v", logOutput, expectedLogOutput)
	}
}
