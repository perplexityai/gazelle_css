package css

import (
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func TestCSSModuleConfigDefaults(t *testing.T) {
	c := config.New()
	NewLanguage().Configure(c, "", nil)

	cfg := c.Exts[languageName].(*cssConfig)
	if cfg.modules.enabled {
		t.Fatal("CSS Modules are enabled by default")
	}
	if got, want := cfg.modules.name, "css"; got != want {
		t.Fatalf("CSS Module rule name = %q, want %q", got, want)
	}
}

func TestCSSModuleConfigInheritsAndOverrides(t *testing.T) {
	c := config.New()
	NewLanguage().Configure(c, "", &rule.File{Directives: []rule.Directive{
		{Key: directiveModuleEnabled, Value: "true"},
		{Key: directiveModuleName, Value: "styles"},
	}})
	NewLanguage().Configure(c, "child", nil)

	inherited := c.Exts[languageName].(*cssConfig)
	if !inherited.modules.enabled {
		t.Fatal("child config did not inherit enabled CSS Modules")
	}
	if got, want := inherited.modules.name, "styles"; got != want {
		t.Fatalf("inherited CSS Module rule name = %q, want %q", got, want)
	}

	NewLanguage().Configure(c, "child/grandchild", &rule.File{Directives: []rule.Directive{
		{Key: directiveModuleEnabled, Value: "false"},
		{Key: directiveModuleName, Value: "contracts"},
	}})
	overridden := c.Exts[languageName].(*cssConfig)
	if overridden.modules.enabled {
		t.Fatal("grandchild config did not disable CSS Modules")
	}
	if got, want := overridden.modules.name, "contracts"; got != want {
		t.Fatalf("overridden CSS Module rule name = %q, want %q", got, want)
	}
}
