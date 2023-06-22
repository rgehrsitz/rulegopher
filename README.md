# Rulegopher

Rulegopher is a rule engine written in Go. It allows you to define rules and facts, and evaluate these facts against the rules. The engine triggers events when certain conditions are met.

## Features

- Define rules with conditions and events.
- Add and remove rules dynamically.
- Evaluate facts against the rules.
- Trigger events based on rule evaluation.

## Structure

The project is structured as follows:

- `cmd/server/main.go`: The main entry point of the application. It sets up the rule engine, fact handler, and API handler, and starts the HTTP server.

- `pkg/engine/engine.go`: Defines the rule engine, which manages the rules and evaluates facts.

- `pkg/rules/rules.go`: Defines the structures for rules, conditions, facts, and events, and provides a method for evaluating a fact against a rule.

- `pkg/facts/facts.go`: Defines a fact handler that uses the rule engine to evaluate facts.

- `api/handler/handler.go`: Defines an API handler that provides HTTP endpoints for adding and removing rules, and evaluating facts.

## Getting Started

To get started with Rulegopher, clone the repository and install the necessary dependencies.

```bash
git clone https://github.com/rgehrsitz/rulegopher.git
cd rulegopher
go mod download
```

## Usage

```bash
go run cmd/server/main.go
```
By default, the server listens on port 8080. You can specify a different port with the -port flag. You can also enable logging with the -logging flag, and specify a JSON file containing initial rules with the -rules flag.

Once the server is running, you can interact with it through the following HTTP endpoints:

- POST /addRule: Adds a new rule. The rule should be provided in the request body as a JSON object.
- GET /removeRule?name=<ruleName>: Removes the rule with the specified name.
- POST /evaluateFact: Evaluates a fact. The fact should be provided in the request body as a JSON object. The response is a list of events triggered by the fact.

## Rule Specification

A rule in Rulegopher is defined as a JSON object with the following properties:

- **name**: A string that uniquely identifies the rule.
- **priority**: An integer that determines the order in which the rules are evaluated. Lower numbers indicate higher priority.
- **conditions**: An object that specifies the conditions under which the rule is triggered. It has two properties:
- **all**: An array of conditions that must all be met for the rule to be triggered.
- **any**: An array of conditions, any of which can be met for the rule to be triggered.
- **event**: An object that specifies the event that is triggered when the rule is met. It has the following properties:
- **eventType**: A string that identifies the type of event.
- **customProperty**: A custom property that can be used to store additional information about the event.
- **facts**: An array of facts that triggered the event. This is populated when the rule is evaluated.
- **values**: An array of values corresponding to the facts that triggered the event. This is populated when the rule is evaluated.

Each condition in the all and any arrays is an object with the following properties:

- **fact**: A string that identifies the fact to be evaluated.
- **operator**: A string that specifies the operator to be used for the evaluation. It can be one of the following: equal, notEqual, greaterThan, greaterThanOrEqual, lessThan, lessThanOrEqual, contains, notContains.
**value**: The value to be compared with the fact.

## Rule Example

Here's an example of a rule that triggers an event when a user's age is greater than or equal to 18:

```json
{
  "name": "AdultUser",
  "priority": 1,
  "conditions": {
    "all": [
      {
        "fact": "age",
        "operator": "greaterThanOrEqual",
        "value": 18
      }
    ]
  },
  "event": {
    "eventType": "UserIsAdult",
    "customProperty": "User has reached adulthood."
  }
}
```

In this example, the rule is named "AdultUser" and has a priority of 1. It has a single condition that checks if the "age" fact is greater than or equal to 18. If this condition is met, it triggers an event of type "UserIsAdult" with a custom property "User has reached adulthood.".

## Acknowledgements

This project was heavily inspired by:

- [json-rules-engine](https://github.com/CacheControl/json-rules-engine)
- [Go-Rules-Engine](https://github.com/Icheka/go-rules-engine)

Thank you to both of those authors for their excellent projects that I learned from.

## License
[MIT](https://choosealicense.com/licenses/mit/)