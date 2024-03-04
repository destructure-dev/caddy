package caddy

import (
	"fmt"
	"strings"
	"sync"
)

var (
	modulesMu sync.RWMutex
	modules   = make(map[string]ModuleInfo)
)

// Module provides configuration for a Caddy module.
type Module interface {
	CaddyModule() ModuleInfo
}

// ModuleInfo describes a Caddy module.
type ModuleInfo struct {
	// ID is a unique, namespaced identifier for the module.
	// Caddy IDs are dot separated and use snake_case by convention.
	// For example, `caddy.storage.file_system` is a module with the
	// name "file_system" in the namespace "caddy.storage".
	ID string

	// New returns a pointer to a new, empty instance of the module's type.
	New func() Module
}

// Name returns the name of the module.
func (m ModuleInfo) Name() string {
	i := strings.LastIndexByte(m.ID, '.')

	if i == -1 {
		return ""
	}

	return m.ID[i+1:]
}

// RegisterModule adds a module to the registry by it's ID.
// Registering a module makes it available for lookup when unmarshalling config.
func RegisterModule(module Module) {
	modulesMu.Lock()
	defer modulesMu.Unlock()

	modules[module.CaddyModule().ID] = module.CaddyModule()
}

// GetModule returns module info for the given ID, if found.
func GetModule(id string) (ModuleInfo, error) {
	modulesMu.RLock()
	defer modulesMu.RUnlock()

	m, ok := modules[id]

	if !ok {
		return ModuleInfo{}, fmt.Errorf("Module not found: %s", id)
	}

	return m, nil
}
