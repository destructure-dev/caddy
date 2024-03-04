package caddy

import (
	"encoding/json"
	"fmt"
)

func init() {
	RegisterModule(FileStorage{})
}

// Storage defines how/where Caddy stores assets (such as TLS certificates).
type Storage struct {
	Module string `json:"module"`
}

// UnmarshalStorage unmarshals a storage module.
func UnmarshalStorage(buf []byte) (Module, error) {
	var s Storage
	if err := json.Unmarshal(buf, &s); err != nil {
		return nil, fmt.Errorf("unmarshal storage: %w", err)
	}

	id := fmt.Sprintf("caddy.storage.%s", s.Module)

	m, err := GetModule(id)

	if err != nil {
		return nil, fmt.Errorf("getting storage module: %w", err)
	}

	sm := m.New()

	if err := json.Unmarshal(buf, sm); err != nil {
		return nil, fmt.Errorf("unmarshal storage module: %w", err)
	}

	return sm, nil
}

// FileStorage configures a local filesystem for certificate storage.
type FileStorage struct {
	Storage
	Root string `json:"root,omitempty"`
}

func NewFileStorage(root string) FileStorage {
	return FileStorage{
		Storage: Storage{
			Module: "file_system",
		},
		Root: root,
	}
}

// CaddyModule implements Module.
func (s FileStorage) CaddyModule() ModuleInfo {
	return ModuleInfo{
		ID: "caddy.storage.file_system",
		New: func() Module {
			return &FileStorage{}
		},
	}
}
