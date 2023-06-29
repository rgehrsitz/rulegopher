package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rgehrsitz/rulegopher/pkg/engine"
	"github.com/rgehrsitz/rulegopher/pkg/facts"
	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

const DefaultPriority = 99

// Handler handles HTTP requests for adding, removing, and evaluating rules.
type Handler struct {
	engine      *engine.Engine
	factHandler *facts.FactHandler
}

// NewHandler creates a new Handler.
func NewHandler(engine *engine.Engine, factHandler *facts.FactHandler) *Handler {
	return &Handler{
		engine:      engine,
		factHandler: factHandler,
	}
}

// AddRule handles HTTP requests for adding a rule.
// It decodes the request body into a Rule and adds it to the engine.
func (h *Handler) AddRule(w http.ResponseWriter, r *http.Request) {
	var rule rules.Rule
	err := json.NewDecoder(r.Body).Decode(&rule)
	if err != nil {
		log.Printf("Failed to decode the request body into a rule: %v", err)
		http.Error(w, "Failed to decode the request body into a rule: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set a default priority if the Priority field is missing
	if rule.Priority == 0 {
		rule.Priority = DefaultPriority
	}

	// Check if the other required fields are missing
	if rule.Name == "" || (len(rule.Conditions.All) == 0 && len(rule.Conditions.Any) == 0) || rule.Event.EventType == "" {
		log.Printf("Missing required fields in the rule")
		http.Error(w, "Missing required fields in the rule", http.StatusBadRequest)
		return
	}

	h.engine.AddRule(rule)
	w.WriteHeader(http.StatusCreated)
}

// RemoveRule handles HTTP requests for removing a rule.
// It gets the rule name from the 'name' query parameter and removes the rule from the engine.
func (h *Handler) RemoveRule(w http.ResponseWriter, r *http.Request) {
	ruleName := r.URL.Query().Get("name")
	if ruleName == "" {
		log.Printf("Missing 'name' query parameter in the request")
		http.Error(w, "Missing 'name' query parameter in the request", http.StatusBadRequest)
		return
	}
	h.engine.RemoveRule(ruleName)
	w.WriteHeader(http.StatusOK)
}

// EvaluateFact handles HTTP requests for evaluating a fact.
// It decodes the request body into a Fact and evaluates it using the fact handler.
func (h *Handler) EvaluateFact(w http.ResponseWriter, r *http.Request) {
	var fact rules.Fact
	err := json.NewDecoder(r.Body).Decode(&fact)
	if err != nil || len(fact) == 0 {
		log.Printf("Failed to decode the request body into a fact: %v", err)
		http.Error(w, "Failed to decode the request body into a fact: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the fact is empty
	if len(fact) == 0 {
		log.Printf("Empty fact in the request body")
		http.Error(w, "Empty fact in the request body", http.StatusBadRequest)
		return
	}

	events := h.factHandler.HandleFact(fact)
	json.NewEncoder(w).Encode(events)
}

// ServeHTTP routes the HTTP request to the appropriate handler function.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/addrule":
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid method. Only POST is allowed.", http.StatusMethodNotAllowed)
			return
		}
		h.AddRule(w, r)
	case "/removerule":
		if r.Method != http.MethodDelete {
			http.Error(w, "Invalid method. Only DELETE is allowed.", http.StatusMethodNotAllowed)
			return
		}
		h.RemoveRule(w, r)
	case "/evaluatefact":
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid method. Only POST is allowed.", http.StatusMethodNotAllowed)
			return
		}
		h.EvaluateFact(w, r)
	default:
		http.Error(w, "Invalid path", http.StatusNotFound)
	}
}
