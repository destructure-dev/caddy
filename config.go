package caddy

import (
	"encoding/json"
	"fmt"
)

// Config defines the Caddy configuration structure.
type Config struct {
	Admin      *AdminConfig               `json:"admin,omitempty"`
	Logging    *Logging                   `json:"logging,omitempty"`
	Storage    Module                     `json:"-"`
	StorageRaw json.RawMessage            `json:"storage,omitempty"`
	Apps       map[string]Module          `json:"-"`
	AppsRaw    map[string]json.RawMessage `json:"apps,omitempty"`
}

// rawConfig is a type alias to allow unmarshalling only raw fields and modules.
type rawConfig Config

// UnmarshalJSON implements json.Unmarshaler.
func (c *Config) UnmarshalJSON(buf []byte) error {
	// decode only raw modules
	if err := json.Unmarshal(buf, (*rawConfig)(c)); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	if len(c.StorageRaw) > 0 {
		s, err := UnmarshalStorage(c.StorageRaw)

		if err != nil {
			return err
		}

		c.Storage = s
	}

	apps, err := UnmarshalApps(c.AppsRaw)

	if err != nil {
		return err
	}

	c.Apps = apps

	return nil
}

// MarshalJSON implements json.Marshaler.
func (c Config) MarshalJSON() ([]byte, error) {
	if c.Storage != nil {
		s, err := json.Marshal(c.Storage)

		if err != nil {
			return nil, err
		}

		c.StorageRaw = s
	}

	rawApps, err := MarshalApps(c.Apps)

	if err != nil {
		return nil, err
	}

	c.AppsRaw = rawApps

	return json.Marshal((rawConfig)(c))
}

// AdminConfig configures Caddy's API endpoint, which is used to manage Caddy while it is running.
type AdminConfig struct {
	ID       string `json:"@id,omitempty"`
	Disabled bool   `json:"disabled"`
	Listen   string `json:"listen,omitempty"`
}

// Logging configures logging within Caddy.
type Logging struct {
	ID string `json:"@id,omitempty"`
	//
}
