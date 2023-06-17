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
