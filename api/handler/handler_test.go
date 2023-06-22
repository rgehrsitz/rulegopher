package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rgehrsitz/rulegopher/pkg/engine"
	"github.com/rgehrsitz/rulegopher/pkg/facts"
	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

func TestAddRule(t *testing.T) {
	e := engine.NewEngine()
	fh := facts.NewFactHandler(e)
	h := NewHandler(e, fh)

	rule := rules.Rule{
		Name:     "TestRule",
		Priority: 1,
		Conditions: rules.Conditions{
			All: []rules.Condition{
				{
					Fact:     "temperature",
					Operator: "greaterThan",
					Value:    30,
				},
			},
		},
		Event: rules.Event{
			EventType:      "alert",
			CustomProperty: "AC turned on",
		},
	}

	ruleJSON, _ := json.Marshal(rule)
	req, _ := http.NewRequest("POST", "/addrule", bytes.NewBuffer(ruleJSON))
	rr := httptest.NewRecorder()
	h.AddRule(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestRemoveRule(t *testing.T) {
	e := engine.NewEngine()
	fh := facts.NewFactHandler(e)
	h := NewHandler(e, fh)

	ruleName := "TestRule"
	req, _ := http.NewRequest("DELETE", "/removerule?name="+ruleName, nil)
	rr := httptest.NewRecorder()
	h.RemoveRule(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestEvaluateFact(t *testing.T) {
	e := engine.NewEngine()
	fh := facts.NewFactHandler(e)
	h := NewHandler(e, fh)

	fact := rules.Fact{
		"temperature": 35,
	}

	factJSON, _ := json.Marshal(fact)
	req, _ := http.NewRequest("POST", "/evaluatefact", bytes.NewBuffer(factJSON))
	rr := httptest.NewRecorder()
	h.EvaluateFact(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHandlerAddRuleWithMissingFields(t *testing.T) {
	// Create a new engine and fact handler
	eng := engine.NewEngine()
	fh := facts.NewFactHandler(eng)

	// Create a new handler with the engine and fact handler
	h := NewHandler(eng, fh)

	// Create a new HTTP request with a rule that has missing fields
	req, err := http.NewRequest("POST", "/addrule", bytes.NewBuffer([]byte(`{
		"Name": "",
		"Conditions": {
			"All": [],
			"Any": []
		},
		"Event": {
			"EventType": ""
		}
	}`)))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the AddRule method
	h.AddRule(rr, req)

	// Check that an HTTP 400 Bad Request error was returned
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestHandlerRemoveRuleWithNonexistentRuleName(t *testing.T) {
	// Create a new engine and fact handler
	eng := engine.NewEngine()
	fh := facts.NewFactHandler(eng)

	// Create a new handler with the engine and fact handler
	h := NewHandler(eng, fh)

	// Create a new HTTP request with a rule name that does not exist
	req, err := http.NewRequest("DELETE", "/removerule?name=NonexistentRule", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the RemoveRule method
	h.RemoveRule(rr, req)

	// Check that an HTTP 200 OK status code was returned
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHandlerEvaluateFactWithInvalidInput(t *testing.T) {
	// Create a new engine and fact handler
	eng := engine.NewEngine()
	fh := facts.NewFactHandler(eng)

	// Create a new handler with the engine and fact handler
	h := NewHandler(eng, fh)

	// Create a new HTTP request with invalid input
	req, err := http.NewRequest("POST", "/evaluatefact", bytes.NewBuffer([]byte(`{"invalid": "input",`))) // malformed JSON
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the EvaluateFact method
	h.EvaluateFact(rr, req)

	// Check that an HTTP 400 Bad Request error was returned
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
