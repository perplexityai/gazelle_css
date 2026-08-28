package css

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
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
	if want := []string{"Button.module.css"}; !reflect.DeepEqual(got.Gen[1].AttrStrings("srcs"), want) {
		t.Errorf("css_module_library srcs = %v, want %v", got.Gen[1].AttrStrings("srcs"), want)
	}
	if got, want := len(got.Imports), len(got.Gen); got != want {
		t.Fatalf("generated %d import payloads for %d rules", got, want)
	}
}

func TestGenerateRulesRequiresCSSModulesToOwnTheirBazelPackage(t *testing.T) {
	cfg := newCSSConfig()
	cfg.modules.enabled = true
	c := &config.Config{Exts: map[string]interface{}{languageName: cfg}}
	var fatalMessage string

	got := newLanguage(func(format string, args ...any) {
		fatalMessage = fmt.Sprintf(format, args...)
	}).GenerateRules(language.GenerateArgs{
		Config:       c,
		Rel:          "components",
		RegularFiles: []string{"nested/Card.module.css"},
	})

	if !strings.Contains(fatalMessage, `CSS Module source "nested/Card.module.css"`) ||
		!strings.Contains(fatalMessage, "nested/BUILD.bazel") {
		t.Fatalf("fatal message %q does not explain how to establish package ownership", fatalMessage)
	}
	if len(got.Gen) != 0 || len(got.Empty) != 0 {
		t.Fatalf("generated rules for a nested CSS Module: Gen=%v Empty=%v", got.Gen, got.Empty)
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

func TestGenerateRulesDeletesStaleCSSModuleRule(t *testing.T) {
	c := &config.Config{Exts: map[string]interface{}{languageName: newCSSConfig()}}
	existing := rule.NewRule(cssModuleKind, "css")
	existing.SetAttr("srcs", []string{"Removed.module.css"})
	got := NewLanguage().GenerateRules(language.GenerateArgs{
		Config: c,
		Rel:    "components",
		File:   &rule.File{Rules: []*rule.Rule{existing}},
	})

	if got, want := len(got.Empty), 1; got != want {
		t.Fatalf("generated %d empty rules, want %d", got, want)
	}
	if got, want := got.Empty[0].Kind(), cssModuleKind; got != want {
		t.Fatalf("empty module rule kind = %q, want %q", got, want)
	}
	if got, want := got.Empty[0].Name(), "css"; got != want {
		t.Fatalf("empty module rule name = %q, want %q", got, want)
	}
}

func TestGenerateRulesSplitsLegacyLibraryWhenModulesAreEnabled(t *testing.T) {
	cfg := newCSSConfig()
	cfg.modules.enabled = true
	c := &config.Config{Exts: map[string]interface{}{languageName: cfg}}
	existing := rule.NewRule(cssLibraryKind, "components")
	existing.SetAttr("srcs", []string{"Button.module.css"})
	got := NewLanguage().GenerateRules(language.GenerateArgs{
		Config:       c,
		Rel:          "components",
		File:         &rule.File{Rules: []*rule.Rule{existing}},
		RegularFiles: []string{"Button.module.css"},
	})

	if count, want := len(got.Gen), 1; count != want || got.Gen[0].Kind() != cssModuleKind {
		t.Fatalf("generated rules = %v, want one %s", got.Gen, cssModuleKind)
	}
	if count, want := len(got.Empty), 1; count != want || got.Empty[0].Kind() != cssLibraryKind {
		t.Fatalf("empty rules = %v, want one %s", got.Empty, cssLibraryKind)
	}
}

func TestGenerateRulesRestoresPlainLibraryWhenModulesAreDisabled(t *testing.T) {
	c := &config.Config{Exts: map[string]interface{}{languageName: newCSSConfig()}}
	existing := rule.NewRule(cssModuleKind, "css")
	existing.SetAttr("srcs", []string{"Button.module.css"})
	got := NewLanguage().GenerateRules(language.GenerateArgs{
		Config:       c,
		Rel:          "components",
		File:         &rule.File{Rules: []*rule.Rule{existing}},
		RegularFiles: []string{"Button.module.css"},
	})

	if count, want := len(got.Gen), 1; count != want || got.Gen[0].Kind() != cssLibraryKind {
		t.Fatalf("generated rules = %v, want one %s", got.Gen, cssLibraryKind)
	}
	if count, want := len(got.Empty), 1; count != want || got.Empty[0].Kind() != cssModuleKind {
		t.Fatalf("empty rules = %v, want one %s", got.Empty, cssModuleKind)
	}
}

func TestGenerateRulesReplacesLegacyLibraryWhenTargetNameIsShared(t *testing.T) {
	tests := []struct {
		name         string
		existingKind string
		configure    func(*config.Config)
	}{
		{name: "canonical library", existingKind: cssLibraryKind},
		{
			name:         "mapped library",
			existingKind: "custom_css_library",
			configure: func(c *config.Config) {
				c.KindMap = map[string]config.MappedKind{
					cssLibraryKind: {
						FromKind: cssLibraryKind,
						KindName: "custom_css_library",
						KindLoad: "//tools:css.bzl",
					},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newCSSConfig()
			cfg.modules.enabled = true
			c := &config.Config{Exts: map[string]interface{}{languageName: cfg}}
			if tt.configure != nil {
				tt.configure(c)
			}
			existing := rule.NewRule(tt.existingKind, "css")
			existing.SetAttr("srcs", []string{"Button.module.css"})
			var fatalMessage string
			got := newLanguage(func(format string, args ...any) {
				fatalMessage = fmt.Sprintf(format, args...)
			}).GenerateRules(language.GenerateArgs{
				Config:       c,
				Rel:          "css",
				File:         &rule.File{Rules: []*rule.Rule{existing}},
				RegularFiles: []string{"Button.module.css"},
			})

			if fatalMessage != "" {
				t.Fatalf("GenerateRules reported a valid transition as a conflict: %s", fatalMessage)
			}
			if count, want := len(got.Gen), 1; count != want || got.Gen[0].Kind() != cssModuleKind {
				t.Fatalf("generated rules = %v, want one %s", got.Gen, cssModuleKind)
			}
			if count, want := len(got.Empty), 1; count != want || got.Empty[0].Kind() != cssLibraryKind {
				t.Fatalf("empty rules = %v, want one %s", got.Empty, cssLibraryKind)
			}
		})
	}
}

func TestGenerateRulesDeletesStaleMappedCSSModuleRule(t *testing.T) {
	c := &config.Config{
		Exts: map[string]interface{}{languageName: newCSSConfig()},
		KindMap: map[string]config.MappedKind{
			cssModuleKind: {
				FromKind: cssModuleKind,
				KindName: "custom_css_module_library",
				KindLoad: "//tools:css.bzl",
			},
		},
	}
	existing := rule.NewRule("custom_css_module_library", "css")
	existing.SetAttr("srcs", []string{"Removed.module.css"})
	got := NewLanguage().GenerateRules(language.GenerateArgs{
		Config: c,
		Rel:    "components",
		File:   &rule.File{Rules: []*rule.Rule{existing}},
	})

	if got, want := len(got.Empty), 1; got != want {
		t.Fatalf("generated %d empty rules for mapped CSS Module, want %d", got, want)
	}
	if got, want := got.Empty[0].Kind(), cssModuleKind; got != want {
		t.Fatalf("empty mapped module rule kind = %q, want canonical %q", got, want)
	}
}

func TestGenerateRulesFailsClosedOnCSSModuleOwnershipConflicts(t *testing.T) {
	otherOwner := rule.NewRule("filegroup", "css")
	misnamedModule := rule.NewRule(cssModuleKind, "legacy_css")
	canonicalModule := rule.NewRule(cssModuleKind, "css")
	mappedModule := rule.NewRule("custom_css_module_library", "css")
	conflictingAfterModule := rule.NewRule("filegroup", "css")
	generatedByOtherLanguage := rule.NewRule("other_generated_rule", "css")

	tests := []struct {
		name           string
		rel            string
		regularFiles   []string
		file           *rule.File
		otherGen       []*rule.Rule
		disableModules bool
		configure      func(*config.Config)
		wantDetail     string
	}{
		{
			name:         "target name owned by another kind",
			rel:          "components",
			regularFiles: []string{"Button.module.css"},
			file:         &rule.File{Rules: []*rule.Rule{otherOwner}},
			wantDetail:   `target "css" is already owned by filegroup`,
		},
		{
			name:         "target generated earlier in the same run",
			rel:          "components",
			regularFiles: []string{"Button.module.css"},
			otherGen:     []*rule.Rule{generatedByOtherLanguage},
			wantDetail:   `target "css" is already owned by other_generated_rule`,
		},
		{
			name:         "canonical module generated by another language",
			rel:          "components",
			regularFiles: []string{"Button.module.css"},
			otherGen:     []*rule.Rule{canonicalModule},
			wantDetail:   `css_module_library(name = "css") is already generated by another language`,
		},
		{
			name:         "mapped module generated by another language",
			rel:          "components",
			regularFiles: []string{"Button.module.css"},
			otherGen:     []*rule.Rule{mappedModule},
			configure: func(c *config.Config) {
				c.KindMap = map[string]config.MappedKind{
					cssModuleKind: {
						FromKind: cssModuleKind,
						KindName: "custom_css_module_library",
						KindLoad: "//tools:css.bzl",
					},
				}
			},
			wantDetail: `custom_css_module_library(name = "css") is already generated by another language`,
		},
		{
			name:         "module rule uses noncanonical name",
			rel:          "components",
			regularFiles: []string{"Button.module.css"},
			file:         &rule.File{Rules: []*rule.Rule{misnamedModule}},
			wantDetail:   `css_module_library(name = "legacy_css") does not use the configured name "css"`,
		},
		{
			name:         "canonical and mapped rules both exist",
			rel:          "components",
			regularFiles: []string{"Button.module.css"},
			file:         &rule.File{Rules: []*rule.Rule{canonicalModule, mappedModule}},
			configure: func(c *config.Config) {
				c.KindMap = map[string]config.MappedKind{
					cssModuleKind: {
						FromKind: cssModuleKind,
						KindName: "custom_css_module_library",
						KindLoad: "//tools:css.bzl",
					},
				}
			},
			wantDetail: `multiple CSS Module rules exist; expected one aggregate target "css"`,
		},
		{
			name: "stale duplicate rules after the last source is removed",
			rel:  "components",
			file: &rule.File{Rules: []*rule.Rule{canonicalModule, mappedModule}},
			configure: func(c *config.Config) {
				c.KindMap = map[string]config.MappedKind{
					cssModuleKind: {
						FromKind: cssModuleKind,
						KindName: "custom_css_module_library",
						KindLoad: "//tools:css.bzl",
					},
				}
			},
			wantDetail: `multiple CSS Module rules exist; expected one aggregate target "css"`,
		},
		{
			name:       "stale rule uses the old configured name",
			rel:        "components",
			file:       &rule.File{Rules: []*rule.Rule{misnamedModule}},
			wantDetail: `css_module_library(name = "legacy_css") does not use the configured name "css"`,
		},
		{
			name:           "disabled modules still validate stale ownership",
			rel:            "components",
			regularFiles:   []string{"Button.module.css"},
			file:           &rule.File{Rules: []*rule.Rule{canonicalModule, conflictingAfterModule}},
			disableModules: true,
			wantDetail:     `target "css" is already owned by filegroup`,
		},
		{
			name:         "another kind collides after valid module rule",
			rel:          "components",
			regularFiles: []string{"Button.module.css"},
			file:         &rule.File{Rules: []*rule.Rule{canonicalModule, conflictingAfterModule}},
			wantDetail:   `target "css" is already owned by filegroup`,
		},
		{
			name:         "generated plain and module targets collide",
			rel:          "css",
			regularFiles: []string{"global.css", "Button.module.css"},
			wantDetail:   `css_library and css_module_library would both generate target "css"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newCSSConfig()
			cfg.modules.enabled = !tt.disableModules
			c := &config.Config{Exts: map[string]interface{}{languageName: cfg}}
			if tt.configure != nil {
				tt.configure(c)
			}
			var fatalMessage string
			lang := newLanguage(func(format string, args ...any) {
				fatalMessage = fmt.Sprintf(format, args...)
			})

			got := lang.GenerateRules(language.GenerateArgs{
				Config:       c,
				Rel:          tt.rel,
				File:         tt.file,
				RegularFiles: tt.regularFiles,
				OtherGen:     tt.otherGen,
			})

			if fatalMessage == "" {
				t.Fatal("GenerateRules did not report the ownership conflict")
			}
			if !strings.Contains(fatalMessage, tt.wantDetail) {
				t.Fatalf("fatal message %q does not contain %q", fatalMessage, tt.wantDetail)
			}
			if !strings.Contains(fatalMessage, "# gazelle:ignore") || !strings.Contains(fatalMessage, "entire BUILD file") {
				t.Fatalf("fatal message lacks whole-BUILD escape hatch: %q", fatalMessage)
			}
			if len(got.Gen) != 0 || len(got.Empty) != 0 {
				t.Fatalf("generated rules after ownership conflict: Gen=%v Empty=%v", got.Gen, got.Empty)
			}
		})
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
