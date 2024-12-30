package types

// Scope is a map of identifiers to their types
type Scope map[string]string

type Environment struct {
	Scopes []Scope
	Cursor int
}

func NewEnvironment() *Environment {
	// Initialize the environment with an empty scope (the top-level scope)
	return &Environment{
		Scopes: []Scope{{}},
	}
}

func (e *Environment) PushScope() {
	// Push a new scope onto the stack
	e.Scopes = append(e.Scopes, Scope{})
	e.Cursor++
}

func (e *Environment) PopScope() {
	// Pop the top scope from the stack
	e.Scopes = e.Scopes[:e.Cursor]
	e.Cursor--
}

func (e *Environment) Set(key, value string) {
	// Set a value in the current scope
	e.Scopes[e.Cursor][key] = value
}

func (e *Environment) Get(key string) (string, bool) {
	// Get a value from the current scope
	value, ok := e.Scopes[e.Cursor][key]
	return value, ok
}
