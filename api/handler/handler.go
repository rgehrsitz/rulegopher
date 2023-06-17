package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rgehrsitz/rulegopher/pkg/engine"
	"github.com/rgehrsitz/rulegopher/pkg/facts"
	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

type Handler struct {
	engine      *engine.Engine
	factHandler *facts.FactHandler
}

func NewHandler(engine *engine.Engine, factHandler *facts.FactHandler) *Handler {
	return &Handler{
		engine:      engine,
		factHandler: factHandler,
	}
}

func (h *Handler) AddRule(w http.ResponseWriter, r *http.Request) {
	var rule rules.Rule
	err := json.NewDecoder(r.Body).Decode(&rule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.engine.AddRule(rule)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) RemoveRule(w http.ResponseWriter, r *http.Request) {
	ruleName := r.URL.Query().Get("name")
	if ruleName == "" {
		http.Error(w, "Missing rule name", http.StatusBadRequest)
		return
	}
	h.engine.RemoveRule(ruleName)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) EvaluateFact(w http.ResponseWriter, r *http.Request) {
	var fact rules.Fact
	err := json.NewDecoder(r.Body).Decode(&fact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	events := h.factHandler.HandleFact(fact)
	json.NewEncoder(w).Encode(events)
}

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
