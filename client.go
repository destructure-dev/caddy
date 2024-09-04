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
	cfg := &Config{}

	if err := c.sendReadConfig(ctx, "config", "", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Exports Caddy's current configuration at the named path.
// The configuration will be unmarshalled into v.
func (c *Client) GetConfigByPath(ctx context.Context, path string, v any) error {
	return c.sendReadConfig(ctx, "config", path, v)
}

// Exports Caddy's current configuration with the given ID.
// The configuration will be unmarshalled into v.
func (c *Client) GetConfigByID(ctx context.Context, id string, v any) error {
	return c.sendReadConfig(ctx, "id", id, v)
}

func (c *Client) sendReadConfig(ctx context.Context, base string, path string, v any) error {
	path, err := url.JoinPath(c.serverAddr, base, path)

	if err != nil {
		return fmt.Errorf("joining URL path: %w", err)
	}

	resp, err := c.sendRequest(ctx, http.MethodGet, path, nil)

	if err != nil {
		return fmt.Errorf("getting config: %w", err)
	}

	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if err := json.Unmarshal(buf, v); err != nil {
		return fmt.Errorf("unmarshalling config JSON: %w", err)
	}

	return nil
}

// Load sets Caddy's configuration, overriding any previous configuration.
func (c *Client) Load(ctx context.Context, cfg *Config) error {
	path, err := url.JoinPath(c.serverAddr, "load")

	if err != nil {
		return fmt.Errorf("joining URL path: %w", err)
	}

	buf, err := json.Marshal(cfg)

	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}

	resp, err := c.sendRequest(ctx, http.MethodPost, path, bytes.NewBuffer(buf))

	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	resp.Body.Close()

	return nil
}

// PostConfig changes Caddy's configuration at the named path using POST semantics.
func (c *Client) PostConfig(ctx context.Context, path string, cfg any) error {
	return c.sendWriteConfig(ctx, http.MethodPost, path, cfg)
}

// PutConfig hanges Caddy's configuration at the named path using PUT semantics.
func (c *Client) PutConfig(ctx context.Context, path string, cfg any) error {
	return c.sendWriteConfig(ctx, http.MethodPut, path, cfg)
}

// PatchConfig hanges Caddy's configuration at the named path using PATCH semantics.
func (c *Client) PatchConfig(ctx context.Context, path string, cfg any) error {
	return c.sendWriteConfig(ctx, http.MethodPatch, path, cfg)
}

func (c *Client) sendWriteConfig(ctx context.Context, method string, path string, cfg any) error {
	path, err := url.JoinPath(c.serverAddr, "config", path)

	if err != nil {
		return fmt.Errorf("joining config path: %w", err)
	}

	buf, err := json.Marshal(cfg)

	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}

	resp, err := c.sendRequest(ctx, method, path, bytes.NewBuffer(buf))

	if err != nil {
		return err
	}

	resp.Body.Close()

	return nil
}

// DeleteConfig removes Caddy's configuration at the named path.
func (c *Client) DeleteConfig(ctx context.Context, path string) error {
	path, err := url.JoinPath(c.serverAddr, "config", path)

	if err != nil {
		return fmt.Errorf("joining URL path: %w", err)
	}

	resp, err := c.sendRequest(ctx, http.MethodDelete, path, nil)

	if err != nil {
		return fmt.Errorf("deleting caddy config: %w", err)
	}

	resp.Body.Close()

	return nil
}

func (c *Client) sendRequest(ctx context.Context, method string, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Add("content-type", "application/json")

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("sending caddy request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)

		resp.Body.Close()

		return nil, fmt.Errorf("response %d error: %s", resp.StatusCode, string(b))
	}

	return resp, nil
}
