package css

type cssConfig struct {
	enabled    bool
	name       string
	visibility []string
}

func newCSSConfig() *cssConfig {
	return &cssConfig{
		enabled:    true,
		visibility: []string{"//visibility:public"},
	}
}

func (c *cssConfig) clone() *cssConfig {
	clone := *c
	clone.visibility = append([]string(nil), c.visibility...)
	return &clone
}
