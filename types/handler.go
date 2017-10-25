package types

import "errors"

// A Handler is a handler specification.
type Handler struct {
	// Name is the unique identifier for a handler.
	Name string `json:"name"`

	// Type is the handler type, i.e. pipe.
	Type string `json:"type"`

	// Mutator is the handler event data mutator.
	Mutator string `json:"mutator,omitempty"`

	// Command is the command to be executed for a pipe handler.
	Command string `json:"command,omitempty"`

	// Timeout is the handler timeout in seconds.
	Timeout int `json:"timeout"`

	Socket HandlerSocket `json:"socket,omitempty"`

	// Handlers is a list of handlers for a handler set.
	Handlers []string `json:"handlers,omitempty"`

	// Filters is a list of filters to evaluate before executing this handler
	Filters []EventFilter `json:"filters,omitempty"`

	// Env is a list of environment variables to use with command execution
	Env []string `json:"env,omitempty"`

	// Environment indicates to which env a handler belongs to
	Environment string `json:"environment"`

	// Organization indicates to which org a handler belongs to
	Organization string `json:"organization"`
}

// HandlerSocket contains configuration for a TCP or UDP handler.
type HandlerSocket struct {
	// Host is the socket peer address.
	Host string `json:"host"`

	// Port is the socket peer port.
	Port int `json:"port"`
}

// Validate returns an error if the handler does not pass validation tests.
func (h *Handler) Validate() error {
	if err := ValidateName(h.Name); err != nil {
		return errors.New("handler name " + err.Error())
	}

	if err := validateHandlerType(h.Type); err != nil {
		return errors.New("handler type " + err.Error())
	}

	if h.Environment == "" {
		return errors.New("environment must be set")
	}

	if h.Organization == "" {
		return errors.New("organization must be set")
	}

	return nil
}

// GetOrg refers to the organization the handler belongs to
func (h *Handler) GetOrg() string {
	return h.Organization
}

// GetEnv refers to the organization the handler belongs to
func (h *Handler) GetEnv() string {
	return h.Environment
}

// FixtureHandler returns a Handler fixture for testing.
func FixtureHandler(name string) *Handler {
	return &Handler{
		Name:         name,
		Type:         "pipe",
		Command:      "command",
		Environment:  "default",
		Organization: "default",
	}
}

// FixtureSocketHandler returns a Handler fixture for testing.
func FixtureSocketHandler(name string, proto string) *Handler {
	handler := FixtureHandler(name)
	handler.Type = proto
	handler.Socket = HandlerSocket{
		Host: "127.0.0.1",
		Port: 3001,
	}
	return handler
}

// FixtureSetHandler returns a Handler fixture for testing.
func FixtureSetHandler(name string, handlers ...string) *Handler {
	handler := FixtureHandler(name)
	handler.Handlers = handlers
	return handler
}
