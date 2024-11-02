package caddy

import (
	"encoding/json"
	"fmt"
	"time"
)

func init() {
	RegisterModule(HTTP{})
}

// The HTTP app provides a robust, production-ready HTTP server.
type HTTP struct {
	ID            string            `json:"@id,omitempty"`
	HTTPPort      int               `json:"http_port,omitempty"`
	HTTPSPort     int               `json:"https_port,omitempty"`
	GracePeriod   time.Duration     `json:"grace_period,omitempty"`
	ShutdownDelay time.Duration     `json:"shutdown_delay,omitempty"`
	Servers       map[string]Server `json:"servers,omitempty"`
}

// CaddyModule implements Module.
func (h HTTP) CaddyModule() ModuleInfo {
	return ModuleInfo{
		ID: "http",
		New: func() Module {
			return &HTTP{}
		},
	}
}

// Server describes an HTTP server.
type Server struct {
	ID                string        `json:"@id,omitempty"`
	Listen            []string      `json:"listen"`
	ReadTimeout       time.Duration `json:"read_timeout,omitempty"`
	ReadHeaderTimeout time.Duration `json:"read_header_timeout,omitempty"`
	WriteTimeout      time.Duration `json:"write_timeout,omitempty"`
	IdleTimeout       time.Duration `json:"idle_timeout,omitempty"`
	KeepAliveInterval time.Duration `json:"keepalive_interval,omitempty"`
	MaxHeaderBytes    int           `json:"max_header_bytes,omitempty"`
	Routes            []Route       `json:"routes,omitempty"`
	AutoHTTPS         *AutoHTTPS    `json:"automatic_https,omitempty"`
}

// Route consists of a set of rules for matching HTTP requests, a list of
// handlers to execute, and optional flow control parameters which customize
// the handling of HTTP requests in a highly flexible and performant manner.
type Route struct {
	ID        string            `json:"@id,omitempty"`
	Group     string            `json:"group,omitempty"`
	Handle    []Module          `json:"-"`
	HandleRaw []json.RawMessage `json:"handle"`
	Match     []MatcherSet      `json:"match"`
}

// rawRoute is a type alias to allow unmarshalling only raw fields and modules.
type rawRoute Route

// UnmarshalJSON implements json.Unmarshaler.
func (r *Route) UnmarshalJSON(buf []byte) error {
	// decode only raw modules
	if err := json.Unmarshal(buf, (*rawRoute)(r)); err != nil {
		return fmt.Errorf("unmarshal route: %w", err)
	}

	r.Handle = make([]Module, len(r.HandleRaw))

	for i, v := range r.HandleRaw {
		h, err := UnmarshalHandler(v)

		if err != nil {
			return err
		}

		r.Handle[i] = h
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
func (r Route) MarshalJSON() ([]byte, error) {
	r.HandleRaw = make([]json.RawMessage, len(r.Handle))

	for i, v := range r.Handle {
		buf, err := json.Marshal(v)

		if err != nil {
			return nil, err
		}

		r.HandleRaw[i] = buf
	}

	return json.Marshal((rawRoute)(r))
}

// MatcherSet is used to qualify a route for a request.
type MatcherSet struct {
	ID   string   `json:"@id,omitempty"`
	Host []string `json:"host,omitempty"`
	Path []string `json:"path,omitempty"`
}

// AutoHTTPS configures or disables automatic HTTPS within a server.
type AutoHTTPS struct {
	ID      string `json:"@id,omitempty"`
	Disable bool   `json:"disable"`
}
