package caddy

import (
	"encoding/json"
	"fmt"
)

// AppMap is a map of app names to app modules.
type AppMap map[string]Module

// UnmarshalApps unmarshals a map of app modules.
func UnmarshalApps(rawApps map[string]json.RawMessage) (AppMap, error) {
	apps := make(AppMap, len(rawApps))

	for k, v := range rawApps {
		m, err := GetModule(k) // apps are not namespaced

		if err != nil {
			return nil, fmt.Errorf("getting app module: %w", err)
		}

		app := m.New()

		if err := json.Unmarshal(v, app); err != nil {
			return nil, fmt.Errorf("unmarshal %s app module: %w", k, err)
		}

		apps[k] = app
	}

	return apps, nil
}

// MarshalApps marshals a map of app modules to raw JSON messages.
func MarshalApps(apps AppMap) (map[string]json.RawMessage, error) {
	rawApps := make(map[string]json.RawMessage)

	for k, v := range apps {
		buf, err := json.Marshal(v)

		if err != nil {
			return nil, err
		}

		rawApps[k] = buf
	}

	return rawApps, nil
}
