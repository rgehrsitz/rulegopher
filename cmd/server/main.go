package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/rgehrsitz/rulegopher/api/handler"
	"github.com/rgehrsitz/rulegopher/api/middleware"
	"github.com/rgehrsitz/rulegopher/pkg/engine"
	"github.com/rgehrsitz/rulegopher/pkg/facts"
)

func main() {
	// Parse the command line arguments for the port
	port := flag.String("port", "8080", "port to listen on")
	flag.Parse()

	// Create the rules engine and fact handler
	rulesEngine := engine.NewEngine()
	factHandler := facts.NewFactHandler(rulesEngine)

	// Create the API handler
	apiHandler := handler.NewHandler(rulesEngine, factHandler)

	// Set up the routes
	// http.HandleFunc("/addRule", apiHandler.AddRule)
	// http.HandleFunc("/removeRule", apiHandler.RemoveRule)
	// http.HandleFunc("/evaluateFact", apiHandler.EvaluateFact)

	// The following are equivalent to the above but with logging middleware
	http.Handle("/addRule", middleware.LoggingMiddleware(http.HandlerFunc(apiHandler.AddRule)))
	http.Handle("/removeRule", middleware.LoggingMiddleware(http.HandlerFunc(apiHandler.RemoveRule)))
	http.Handle("/evaluateFact", middleware.LoggingMiddleware(http.HandlerFunc(apiHandler.EvaluateFact)))

	// Start the server
	fmt.Printf("Starting server on port %s\n", *port)
	http.ListenAndServe(":"+*port, nil)
}
