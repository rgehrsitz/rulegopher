package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rgehrsitz/rulegopher/api/handler"
	"github.com/rgehrsitz/rulegopher/pkg/engine"
	"github.com/rgehrsitz/rulegopher/pkg/facts"
	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

func TestIntegration(t *testing.T) {
	e := engine.NewEngine()
	fh := facts.NewFactHandler(e)
	h := handler.NewHandler(e, fh)

	server := httptest.NewServer(h)
	defer server.Close()

	// Define a rule
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

	// Send a request to add the rule
	ruleJSON, _ := json.Marshal(rule)
	resp, err := http.Post(server.URL+"/addrule", "application/json", bytes.NewBuffer(ruleJSON))
	if err != nil {
		t.Fatalf("Failed to send request to add rule: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", resp.StatusCode, http.StatusCreated)
	}

	// Define a fact that satisfies the rule
	fact := rules.Fact{
		"temperature": 35,
	}

	// Send a request to evaluate the fact
	factJSON, _ := json.Marshal(fact)
	resp, err = http.Post(server.URL+"/evaluatefact", "application/json", bytes.NewBuffer(factJSON))
	if err != nil {
		t.Fatalf("Failed to send request to evaluate fact: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	// Check that the correct event was returned
	body, _ := ioutil.ReadAll(resp.Body)
	var events []rules.Event
	json.Unmarshal(body, &events)
	if len(events) != 1 || events[0].EventType != "alert" {
		t.Errorf("HandleFact returned incorrect events: got %v", events)
	}
}
