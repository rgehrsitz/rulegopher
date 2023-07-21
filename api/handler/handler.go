package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rgehrsitz/rulegopher/pkg/engine"
	"github.com/rgehrsitz/rulegopher/pkg/facts"
	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

// Handler is a struct that contains an engine and a factHandler.
// @property engine - The `engine` property is a pointer to an instance of the `Engine` struct. It is
// likely used to interact with the main engine or core functionality of the application.
// @property factHandler - The `factHandler` property is an instance of the `FactHandler` struct from
// the `facts` package. It is responsible for handling facts related to the application's logic or
// data.
type Handler struct {
	engine      *engine.Engine
	factHandler *facts.FactHandler
}

// NewHandler returns a new instance of the Handler struct with the provided engine and
// factHandler.
func NewHandler(engine *engine.Engine, factHandler *facts.FactHandler) *Handler {
	return &Handler{
		engine:      engine,
		factHandler: factHandler,
	}
}

// AddRule is a method of the `Handler` struct. It is responsible for adding a new rule
// to the engine.
func (h *Handler) AddRule(w http.ResponseWriter, r *http.Request) {
	var rule rules.Rule
	err := json.NewDecoder(r.Body).Decode(&rule)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Set a default priority if the Priority field is missing
	if rule.Priority == 0 {
		rule.Priority = 99
	}

	// Check if the other required fields are missing
	if rule.Name == "" || (len(rule.Conditions.All) == 0 && len(rule.Conditions.Any) == 0) || rule.Event.EventType == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	h.engine.AddRule(rule)
	w.WriteHeader(http.StatusCreated)
}

// RemoveRule is a method of the `Handler` struct. It is responsible for removing a rule
// from the engine based on the provided rule name.
func (h *Handler) RemoveRule(w http.ResponseWriter, r *http.Request) {
	ruleName := r.URL.Query().Get("name")
	if ruleName == "" {
		http.Error(w, "Missing rule name", http.StatusBadRequest)
		return
	}
	h.engine.RemoveRule(ruleName)
	w.WriteHeader(http.StatusOK)
}

// EvaluateFact is a method of the `Handler` struct. It is responsible for evaluating a
// fact by decoding the fact data from the request body, handling the fact using the `factHandler`
// instance, and encoding the resulting events as a JSON response.
func (h *Handler) EvaluateFact(w http.ResponseWriter, r *http.Request) {
	var fact rules.Fact
	err := json.NewDecoder(r.Body).Decode(&fact)
	if err != nil || len(fact) == 0 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	events, err := h.factHandler.HandleFact(fact)
	if err != nil {
		http.Error(w, "Error evaluating fact", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(events)
}

// ServeHTTP` is a method of the `Handler` struct that implements the `http.Handler`
// interface. It is responsible for handling incoming HTTP requests and routing them to the appropriate
// methods based on the URL path.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/addrule":
		h.AddRule(w, r)
	case "/removerule":
		h.RemoveRule(w, r)
	case "/evaluatefact":
		h.EvaluateFact(w, r)
	default:
		http.NotFound(w, r)
	}
}
