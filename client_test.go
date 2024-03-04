package caddy_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alecthomas/assert/v2"
	"go.destructure.dev/caddy"
)

func TestGetConfig(t *testing.T) {
	cfg := caddy.Config{
		Storage: caddy.NewFileStorage("/var/caddy"),
	}
	buf, err := json.Marshal(cfg)
	assert.NoError(t, err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/config", r.URL.Path)

		w.Header().Add("content-type", "application/json")
		w.Write(buf)
	}))
	defer ts.Close()

	c := caddy.NewClient(ts.URL)

	res, err := c.GetConfig(context.Background())

	assert.NoError(t, err)

	assert.Equal(t, "/var/caddy", res.Storage.(*caddy.FileStorage).Root)
}
