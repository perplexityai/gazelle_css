package css

type cssConfig struct {
	enabled    bool
	name       string
	visibility []string
	modules    cssModuleConfig
}

type cssModuleConfig struct {
	enabled bool
	name    string
}

func newCSSConfig() *cssConfig {
	return &cssConfig{
		enabled:    true,
		visibility: []string{"//visibility:public"},
		modules: cssModuleConfig{
			enabled: false,
			name:    "css",
		},
	}
}

func (c *cssConfig) clone() *cssConfig {
	clone := *c
	clone.visibility = append([]string(nil), c.visibility...)
	return &clone
}
