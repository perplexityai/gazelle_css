package css

import (
	"reflect"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
)

func TestGenerateRules(t *testing.T) {
	c := &config.Config{Exts: map[string]interface{}{languageName: newCSSConfig()}}
	got := NewLanguage().GenerateRules(language.GenerateArgs{
		Config:       c,
		Rel:          "styles",
		RegularFiles: []string{"reset.css", "theme.CSS", "README.md"},
	})
	if len(got.Gen) != 1 {
		t.Fatalf("got %d rules, want 1", len(got.Gen))
	}
	if got.Gen[0].Name() != "styles" {
		t.Errorf("name = %q, want styles", got.Gen[0].Name())
	}
	if want := []string{"reset.css", "theme.CSS"}; !reflect.DeepEqual(got.Gen[0].AttrStrings("srcs"), want) {
		t.Errorf("srcs = %v, want %v", got.Gen[0].AttrStrings("srcs"), want)
	}
}

func TestGenerateRulesSeparatesCSSModules(t *testing.T) {
	cfg := newCSSConfig()
	cfg.modules.enabled = true
	c := &config.Config{Exts: map[string]interface{}{languageName: cfg}}
	got := NewLanguage().GenerateRules(language.GenerateArgs{
		Config: c,
		Rel:    "components",
		RegularFiles: []string{
			"global.css",
			"Button.module.css",
			"nested/Card.module.css",
			"README.md",
		},
	})

	if got, want := len(got.Gen), 2; got != want {
		t.Fatalf("generated %d rules, want %d", got, want)
	}
	if got, want := got.Gen[0].Kind(), cssLibraryKind; got != want {
		t.Fatalf("first rule kind = %q, want %q", got, want)
	}
	if want := []string{"global.css"}; !reflect.DeepEqual(got.Gen[0].AttrStrings("srcs"), want) {
		t.Errorf("css_library srcs = %v, want %v", got.Gen[0].AttrStrings("srcs"), want)
	}
	if got, want := got.Gen[1].Kind(), cssModuleKind; got != want {
		t.Fatalf("second rule kind = %q, want %q", got, want)
	}
	if got, want := got.Gen[1].Name(), "css"; got != want {
		t.Errorf("css_module_library name = %q, want %q", got, want)
	}
	if want := []string{"Button.module.css", "nested/Card.module.css"}; !reflect.DeepEqual(got.Gen[1].AttrStrings("srcs"), want) {
		t.Errorf("css_module_library srcs = %v, want %v", got.Gen[1].AttrStrings("srcs"), want)
	}
	if got, want := len(got.Imports), len(got.Gen); got != want {
		t.Fatalf("generated %d import payloads for %d rules", got, want)
	}
}

func TestGenerateRulesTreatsModulesAsPlainCSSWhenDisabled(t *testing.T) {
	cfg := newCSSConfig()
	c := &config.Config{Exts: map[string]interface{}{languageName: cfg}}
	got := NewLanguage().GenerateRules(language.GenerateArgs{
		Config:       c,
		Rel:          "components",
		RegularFiles: []string{"Button.module.css"},
	})

	if got, want := len(got.Gen), 1; got != want {
		t.Fatalf("generated %d rules, want %d", got, want)
	}
	if got, want := got.Gen[0].Kind(), cssLibraryKind; got != want {
		t.Fatalf("generated rule kind = %q, want %q", got, want)
	}
	if want := []string{"Button.module.css"}; !reflect.DeepEqual(got.Gen[0].AttrStrings("srcs"), want) {
		t.Fatalf("css_library srcs = %v, want %v", got.Gen[0].AttrStrings("srcs"), want)
	}
}

func TestGenerateRulesDisabled(t *testing.T) {
	c := &config.Config{Exts: map[string]interface{}{
		languageName: &cssConfig{enabled: false},
	}}
	got := NewLanguage().GenerateRules(language.GenerateArgs{Config: c, RegularFiles: []string{"app.css"}})
	if len(got.Gen) != 0 {
		t.Fatalf("got %d rules, want 0", len(got.Gen))
	}
}
