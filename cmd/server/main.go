package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/rgehrsitz/rulegopher/api/handler"
	"github.com/rgehrsitz/rulegopher/api/middleware"
	"github.com/rgehrsitz/rulegopher/pkg/engine"
	"github.com/rgehrsitz/rulegopher/pkg/facts"
	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

func main() {
	// Parse the command line arguments for the port and rules file
	port := flag.String("port", "8080", "port to listen on")
	logging := flag.Bool("logging", false, "enable or disable logging")
	rulesFile := flag.String("rules", "", "JSON file containing the rules")
	reportFacts := flag.Bool("reportFacts", false, "whether to report the facts that caused the event to trigger")
	reportRuleName := flag.Bool("reportRuleName", true, "whether to report the name of the rule that was triggered")

	flag.Parse()

	// Create the rules engine and fact handler
	rulesEngine := engine.NewEngine()
	rulesEngine.ReportFacts = *reportFacts
	rulesEngine.ReportRuleName = *reportRuleName
	factHandler := facts.NewFactHandler(rulesEngine)

	// If a rules file is provided, load the rules from the file
	if *rulesFile != "" {
		file, err := os.Open(*rulesFile)
		if err != nil {
			fmt.Println("Failed to open rules file:", err)
			return
		}
		defer file.Close()

		var rules []rules.Rule
		if err := json.NewDecoder(file).Decode(&rules); err != nil {
			fmt.Println("Failed to decode rules file:", err)
			return
		}

		for _, rule := range rules {
			if err := rulesEngine.AddRule(rule); err != nil {
				fmt.Println("Failed to add rule:", err)
				return
			}
		}
	}

	// Create the API handler
	apiHandler := handler.NewHandler(rulesEngine, factHandler)

	// Set up the routes
	if *logging {
		http.Handle("/addRule", middleware.LoggingMiddleware(http.HandlerFunc(apiHandler.AddRule)))
		http.Handle("/removeRule", middleware.LoggingMiddleware(http.HandlerFunc(apiHandler.RemoveRule)))
		http.Handle("/evaluateFact", middleware.LoggingMiddleware(http.HandlerFunc(apiHandler.EvaluateFact)))
	} else {
		http.Handle("/addRule", http.HandlerFunc(apiHandler.AddRule))
		http.Handle("/removeRule", http.HandlerFunc(apiHandler.RemoveRule))
		http.Handle("/evaluateFact", http.HandlerFunc(apiHandler.EvaluateFact))
	}

	// Start the server
	fmt.Printf("Starting server on port %s\n", *port)
	http.ListenAndServe(":"+*port, nil)
	fmt.Println("Server started")
}
