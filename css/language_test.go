package css

import (
	"reflect"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
)

func TestGenerateRules(t *testing.T) {
	c := &config.Config{Exts: map[string]interface{}{
		languageName: &cssConfig{enabled: true, visibility: []string{"//visibility:public"}},
	}}
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

func TestGenerateRulesDisabled(t *testing.T) {
	c := &config.Config{Exts: map[string]interface{}{
		languageName: &cssConfig{enabled: false},
	}}
	got := NewLanguage().GenerateRules(language.GenerateArgs{Config: c, RegularFiles: []string{"app.css"}})
	if len(got.Gen) != 0 {
		t.Fatalf("got %d rules, want 0", len(got.Gen))
	}
}
