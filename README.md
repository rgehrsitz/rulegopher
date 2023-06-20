# Rulegopher

Rulegopher is a simple rule engine written in Go.

It consists of three main files:

* `facts.go` defines the `Fact` struct, which represents a fact that can be evaluated by the rule engine.
* `rules.go` defines the `Rule` struct, which represents a rule that consists of a set of conditions and an event.
* `handler.go` defines the `Handler` struct, which implements an HTTP handler that provides endpoints for adding, removing, and evaluating rules.

The project also includes a `middleware.go` file that defines a middleware function that logs the start and end time of each HTTP request.

Rulegopher can be used to implement simple rule-based systems. For example, you could use it to implement a system that triggers an event when a certain condition is met, such as when a user enters a certain value into a form.

## Key features

* Simple and easy-to-use rule engine
* Written in Go, a modern and efficient programming language
* Open source, so you can modify and extend it as needed
