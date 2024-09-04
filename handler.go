package caddy

import (
	"encoding/json"
	"fmt"
	"time"
)

func init() {
	RegisterModule(ReverseProxy{})
}

// Handler is a middleware which handles requests.
type Handler struct {
	ID      string `json:"@id,omitempty"`
	Handler string `json:"handler"`
}

// UnmarshalHandler unmarshals a handler module.
func UnmarshalHandler(buf []byte) (Module, error) {
	var h Handler
	if err := json.Unmarshal(buf, &h); err != nil {
		return nil, fmt.Errorf("unmarshal handler: %w", err)
	}

	id := fmt.Sprintf("http.handlers.%s", h.Handler)

	m, err := GetModule(id)

	if err != nil {
		return nil, fmt.Errorf("getting handler module: %w", err)
	}

	hm := m.New()

	if err := json.Unmarshal(buf, hm); err != nil {
		return nil, fmt.Errorf("unmarshal handler module: %w", err)
	}

	return hm, nil
}

// ReverseProxy configures a highly configurable and production-ready reverse proxy.
type ReverseProxy struct {
	Handler
	Upstreams      []Upstream    `json:"upstreams,omitempty"`
	HealthChecks   *HealthChecks `json:"health_checks,omitempty"`
	TrustedProxies []string      `json:"trusted_proxies,omitempty"`
}

func NewReverseProxy() ReverseProxy {
	return ReverseProxy{
		Handler: Handler{
			Handler: "reverse_proxy",
		},
	}
}

// CaddyModule implements Module.
func (s ReverseProxy) CaddyModule() ModuleInfo {
	return ModuleInfo{
		ID: "http.handlers.reverse_proxy",
		New: func() Module {
			return &ReverseProxy{}
		},
	}
}

type Upstream struct {
	ID          string `json:"@id,omitempty"`
	Dial        string `json:"dial"`
	MaxRequests int    `json:"max_requests,omitempty"`
}

type HealthChecks struct {
	ID     string             `json:"@id,omitempty"`
	Active *ActiveHealthCheck `json:"active,omitempty"`
}

type ActiveHealthCheck struct {
	ID       string        `json:"@id,omitempty"`
	URI      string        `json:"uri,omitempty"`
	Interval time.Duration `json:"interval,omitempty"`
	Timeout  time.Duration `json:"timeout,omitempty"`
}
