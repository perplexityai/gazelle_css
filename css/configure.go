package css

import (
	"flag"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func (l *cssLang) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {}

func (l *cssLang) CheckFlags(fs *flag.FlagSet, c *config.Config) error { return nil }

func (l *cssLang) KnownDirectives() []string {
	return []string{"css_extension", "css_library_name", "css_visibility"}
}

func (l *cssLang) Configure(c *config.Config, rel string, f *rule.File) {
	cfg := newCSSConfig()
	if parent, ok := c.Exts[languageName].(*cssConfig); ok {
		cfg = parent.clone()
	}
	if f != nil {
		for _, d := range f.Directives {
			switch d.Key {
			case "css_extension":
				cfg.enabled = strings.TrimSpace(d.Value) != "disabled"
			case "css_library_name":
				cfg.name = strings.TrimSpace(d.Value)
			case "css_visibility":
				cfg.visibility = strings.Fields(d.Value)
			}
		}
	}
	c.Exts[languageName] = cfg
}
