package caddy

func init() {
	RegisterModule(TLS{})
}

// TLS configures TLS facilities including certificate loading and management,
// client auth, and more.
type TLS struct {
	// ...
}

// CaddyModule implements Module.
func (t TLS) CaddyModule() ModuleInfo {
	return ModuleInfo{
		ID: "tls",
		New: func() Module {
			return &TLS{}
		},
	}
}
