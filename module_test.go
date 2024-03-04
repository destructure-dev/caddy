package caddy_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"go.destructure.dev/caddy"
)

func init() {
	caddy.RegisterModule(TestModule{})
}

type TestModule struct {
	//
}

// CaddyModule implements caddy.Module.
func (s TestModule) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "caddy.testing.fakemod",
		New: func() caddy.Module {
			return &TestModule{}
		},
	}
}

func TestGetModule(t *testing.T) {
	m, err := caddy.GetModule("caddy.testing.fakemod")

	assert.NoError(t, err)

	assert.Equal(t, m.ID, "caddy.testing.fakemod")
}

func TestModuleName(t *testing.T) {
	m := TestModule{}

	name := m.CaddyModule().Name()

	assert.Equal(t, "fakemod", name)
}
