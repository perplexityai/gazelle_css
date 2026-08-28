// Package css implements a small Gazelle extension for CSS source packages.
package css

import (
	"flag"
	"path"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	languageName = "css"
	ruleKind     = "css_library"
)

type cssConfig struct {
	enabled    bool
	name       string
	visibility []string
}

type cssLang struct{}

// NewLanguage creates the CSS Gazelle language extension.
func NewLanguage() language.Language { return &cssLang{} }

func (l *cssLang) Name() string { return languageName }

func (l *cssLang) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {}
func (l *cssLang) CheckFlags(fs *flag.FlagSet, c *config.Config) error { return nil }
func (l *cssLang) KnownDirectives() []string {
	return []string{"css_extension", "css_library_name", "css_visibility"}
}

func (l *cssLang) Configure(c *config.Config, rel string, f *rule.File) {
	cfg := &cssConfig{enabled: true, visibility: []string{"//visibility:public"}}
	if parent, ok := c.Exts[languageName].(*cssConfig); ok {
		cfg.enabled = parent.enabled
		cfg.name = parent.name
		cfg.visibility = append([]string(nil), parent.visibility...)
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

func (l *cssLang) Kinds() map[string]rule.KindInfo {
	return map[string]rule.KindInfo{
		ruleKind: {
			NonEmptyAttrs:  map[string]bool{"name": true, "srcs": true},
			MergeableAttrs: map[string]bool{"srcs": true},
			ResolveAttrs:   map[string]bool{"deps": true},
		},
	}
}

func (l *cssLang) Loads() []rule.LoadInfo {
	return []rule.LoadInfo{{Name: "@gazelle_css//css:defs.bzl", Symbols: []string{ruleKind}}}
}

func (l *cssLang) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	cfg, ok := args.Config.Exts[languageName].(*cssConfig)
	if !ok || !cfg.enabled {
		return language.GenerateResult{}
	}
	var srcs []string
	for _, name := range args.RegularFiles {
		if strings.HasSuffix(strings.ToLower(name), ".css") {
			srcs = append(srcs, name)
		}
	}
	if len(srcs) == 0 {
		return language.GenerateResult{}
	}
	sort.Strings(srcs)
	name := cfg.name
	if name == "" {
		name = path.Base(args.Rel)
		if name == "." || name == "" {
			name = "css"
		}
	}
	r := rule.NewRule(ruleKind, name)
	r.SetAttr("srcs", srcs)
	r.SetAttr("visibility", cfg.visibility)
	// Gazelle requires one import-data entry per generated rule. CSS import
	// resolution is intentionally deferred, so the corresponding value is nil.
	return language.GenerateResult{
		Gen:     []*rule.Rule{r},
		Imports: []interface{}{nil},
	}
}

func (l *cssLang) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	if r.Kind() != ruleKind {
		return nil
	}
	var imports []resolve.ImportSpec
	for _, src := range r.AttrStrings("srcs") {
		imports = append(imports, resolve.ImportSpec{Lang: languageName, Imp: path.Join(f.Pkg, src)})
	}
	return imports
}

func (l *cssLang) Resolve(c *config.Config, ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, imports interface{}, from label.Label) {
}

func (l *cssLang) Embeds(r *rule.Rule, from label.Label) []label.Label { return nil }
func (l *cssLang) Fix(c *config.Config, f *rule.File)                    {}
