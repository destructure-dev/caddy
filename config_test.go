package caddy_test

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
	"go.destructure.dev/caddy"
)

func TestUnmarshalStorage(t *testing.T) {
	f, err := os.Open("testdata/storage.json")

	assert.NoError(t, err)

	buf, err := io.ReadAll(f)

	assert.NoError(t, err)

	c := caddy.Config{}

	err = json.Unmarshal(buf, &c)

	assert.NoError(t, err)

	s, ok := c.Storage.(*caddy.FileStorage)

	assert.True(t, ok)

	assert.Equal(t, "file_system", s.Module)
	assert.Equal(t, "/var/caddy", s.Root)
}

func TestMarshalStorage(t *testing.T) {
	c := caddy.Config{
		Storage: caddy.NewFileStorage("/var/caddy"),
	}

	buf, err := json.Marshal(c)

	assert.NoError(t, err)

	assert.Equal(t, []byte(`{"storage":{"module":"file_system","root":"/var/caddy"}}`), buf)
}

func TestUnmarshalHTTP(t *testing.T) {
	f, err := os.Open("testdata/http.json")

	assert.NoError(t, err)

	buf, err := io.ReadAll(f)

	assert.NoError(t, err)

	c := caddy.Config{}

	err = json.Unmarshal(buf, &c)

	assert.NoError(t, err)

	h, ok := c.Apps["http"].(*caddy.HTTP)

	assert.True(t, ok)

	assert.Equal(t, 80, h.HTTPPort)
}

func TestMarshalApps(t *testing.T) {
	c := caddy.Config{
		Apps: map[string]caddy.Module{
			"http": caddy.HTTP{
				HTTPPort: 80,
			},
		},
	}

	buf, err := json.Marshal(c)

	assert.NoError(t, err)

	assert.Equal(t, []byte(`{"apps":{"http":{"http_port":80}}}`), buf)
}

func TestUnmarshalHandler(t *testing.T) {
	f, err := os.Open("testdata/handler.json")

	assert.NoError(t, err)

	buf, err := io.ReadAll(f)

	assert.NoError(t, err)

	c := caddy.Config{}

	err = json.Unmarshal(buf, &c)

	assert.NoError(t, err)

	http, ok := c.Apps["http"].(*caddy.HTTP)

	assert.True(t, ok)

	h, ok := http.Servers["web"].Routes[0].Handle[0].(*caddy.ReverseProxy)

	assert.True(t, ok)

	assert.Equal(t, 1, len(h.Upstreams))
}
