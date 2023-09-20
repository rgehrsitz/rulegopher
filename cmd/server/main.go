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

	// The code block is using the `flag` package in Go to define and parse command-line
	// flags.
	port := flag.String("port", "8080", "port to listen on")
	logging := flag.Bool("logging", false, "enable or disable logging")
	rulesFile := flag.String("rules", "", "JSON file containing the rules")
	reportFacts := flag.Bool("reportFacts", false, "whether to report the facts that caused the event to trigger")
	reportRuleName := flag.Bool("reportRuleName", true, "whether to report the name of the rule that was triggered")
	unmatchedFactBehavior := flag.String("unmatchedFactBehavior", "Ignore", "behavior for unmatched facts: Ignore, Log, or Error")

	flag.Parse()

	// This block of code is creating a new instance of the rules engine and fact handler.
	rulesEngine := engine.NewEngine()
	rulesEngine.ReportFacts = *reportFacts
	rulesEngine.ReportRuleName = *reportRuleName
	rulesEngine.UnmatchedFactBehavior = *unmatchedFactBehavior
	factHandler := facts.NewFactHandler(rulesEngine)

	// This block of code is responsible for reading and decoding the rules from a JSON file, and then
	// adding those rules to the rules engine.
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

	// The line `apiHandler := handler.NewHandler(rulesEngine, factHandler)` is creating a new instance of
	// the `Handler` struct from the `handler` package. It is passing the `rulesEngine` and `factHandler`
	// as arguments to the `NewHandler` function, which initializes the `Handler` struct with these
	// dependencies. This `apiHandler` instance will be used to handle incoming API requests.
	apiHandler := handler.NewHandler(rulesEngine, factHandler)

	// This block of code is responsible for setting up the HTTP handlers for different API endpoints based
	// on the value of the `logging` flag.
	if *logging {
		http.Handle("/addRule", middleware.LoggingMiddleware(http.HandlerFunc(apiHandler.AddRule)))
		http.Handle("/removeRule", middleware.LoggingMiddleware(http.HandlerFunc(apiHandler.RemoveRule)))
		http.Handle("/evaluateFact", middleware.LoggingMiddleware(http.HandlerFunc(apiHandler.EvaluateFact)))
	} else {
		http.Handle("/addRule", http.HandlerFunc(apiHandler.AddRule))
		http.Handle("/removeRule", http.HandlerFunc(apiHandler.RemoveRule))
		http.Handle("/evaluateFact", http.HandlerFunc(apiHandler.EvaluateFact))
	}

	// This code block is responsible for starting the HTTP server and listening for incoming requests on
	// the specified port.
	fmt.Printf("Starting server on port %s\n", *port)
	http.ListenAndServe(":"+*port, nil)
	fmt.Println("Server started")
}
