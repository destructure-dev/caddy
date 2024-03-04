package caddy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
)

// DefaultServerAddr is the default address for a Caddy API server.
const DefaultServerAddr = "http://localhost:2019"

// Client provides an API client for the Caddy server.
type Client struct {
	serverAddr string
	httpClient *http.Client
}

// NewClient creates a new Caddy API client.
func NewClient(address string) *Client {
	if address == "" {
		address = DefaultServerAddr
	}

	return &Client{
		serverAddr: address,
		httpClient: new(http.Client),
	}
}

// NewSocketClient creates a new Caddy API client using a unix socket.
func NewSocketClient(sockAddress string) *Client {
	httpClient := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "unix", sockAddress)
			},
		},
	}

	return &Client{
		serverAddr: "http://127.0.0.1",
		httpClient: &httpClient,
	}
}

// Exports Caddy's current configuration.
func (c *Client) GetConfig(ctx context.Context) (*Config, error) {
	p, err := url.JoinPath(c.serverAddr, "config")

	if err != nil {
		return nil, fmt.Errorf("joining URL path: %w", err)
	}

	resp, err := c.httpClient.Get(p)

	if err != nil {
		return nil, fmt.Errorf("getting config: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	buf, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	cfg := &Config{}

	if err := json.Unmarshal(buf, cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling config JSON: %w", err)
	}

	return cfg, nil
}

// Load sets Caddy's configuration, overriding any previous configuration.
func (c *Client) Load(ctx context.Context, cfg *Config) error {
	p, err := url.JoinPath(c.serverAddr, "load")

	if err != nil {
		return fmt.Errorf("joining URL path: %w", err)
	}

	buf, err := json.Marshal(cfg)

	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}

	resp, err := c.httpClient.Post(p, "application/json", bytes.NewBuffer(buf))

	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("response %d error: %s", resp.StatusCode, string(b))
	}

	return nil
}

// PostConfig changes Caddy's configuration at the named path using POST semantics.
func (c *Client) PostConfig(path string, cfg any) error {
	return c.sendConfig(http.MethodPost, path, cfg)
}

// PutConfig hanges Caddy's configuration at the named path using PUT semantics.
func (c *Client) PutConfig(path string, cfg any) error {
	return c.sendConfig(http.MethodPut, path, cfg)
}

// PatchConfig hanges Caddy's configuration at the named path using PATCH semantics.
func (c *Client) PatchConfig(path string, cfg any) error {
	return c.sendConfig(http.MethodPatch, path, cfg)
}

func (c *Client) sendConfig(method string, path string, cfg any) error {
	p, err := url.JoinPath(c.serverAddr, "config", path)

	if err != nil {
		return fmt.Errorf("joining URL path: %w", err)
	}

	buf, err := json.Marshal(cfg)

	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}

	req, err := http.NewRequest(method, p, bytes.NewBuffer(buf))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Add("content-type", "application/json")

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return fmt.Errorf("changing config: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("response %d error: %s", resp.StatusCode, string(b))
	}

	return nil
}

// DeleteConfig removes Caddy's configuration at the named path.
func (c *Client) DeleteConfig(path string) error {
	p, err := url.JoinPath(c.serverAddr, "config", path)

	if err != nil {
		return fmt.Errorf("joining URL path: %w", err)
	}

	req, err := http.NewRequest(http.MethodDelete, p, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return fmt.Errorf("deleting config: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("response %d error: %s", resp.StatusCode, string(b))
	}

	return nil
}
