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

func TestLoggingMiddlewareTypes(t *testing.T) {
	// Create a dummy HTTP handler that returns an HTTP 200 OK status code
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	// Wrap the dummy HTTP handler with the LoggingMiddleware function
	handler := LoggingMiddleware(next)

	// Define the different types of requests to test
	requests := []struct {
		method string
		url    string
	}{
		{"GET", "/"},
		{"POST", "/"},
		{"PUT", "/"},
		{"DELETE", "/"},
	}

	// Test the LoggingMiddleware function with each request
	for _, request := range requests {
		req, err := http.NewRequest(request.method, request.url, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	}
}

func TestLoggingMiddlewareLogOutput(t *testing.T) {
	// Create a dummy HTTP handler that returns an HTTP 200 OK status code
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	// Wrap the dummy HTTP handler with the LoggingMiddleware function
	handler := LoggingMiddleware(next)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Redirect the standard output to a buffer
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// Call the LoggingMiddleware function
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check that the log output contains the method and URL of the request
	expected := "GET /"
	if !strings.Contains(buf.String(), expected) {
		t.Errorf("Expected log output to contain %q, but it was %q", expected, buf.String())
	}
}

func TestLoggingMiddlewareWithDifferentPaths(t *testing.T) {
	// Create a dummy HTTP handler that returns an HTTP 200 OK status code
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	// Wrap the dummy HTTP handler with the LoggingMiddleware function
	handler := LoggingMiddleware(next)

	// Define the different paths to test
	paths := []string{"/", "/path1", "/path2", "/path3"}

	// Test the LoggingMiddleware function with each path
	for _, path := range paths {
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Redirect the standard output to a buffer
		var buf bytes.Buffer
		log.SetOutput(&buf)

		// Call the LoggingMiddleware function
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		// Check that the log output contains the path
		expected := path
		if !strings.Contains(buf.String(), expected) {
			t.Errorf("Expected log output to contain %q, but it was %q", expected, buf.String())
		}
	}
}
