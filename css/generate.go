package css

import (
	"path"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func (l *cssLang) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	cfg, ok := args.Config.Exts[languageName].(*cssConfig)
	if !ok || !cfg.enabled {
		return language.GenerateResult{}
	}
	var cssSrcs, moduleSrcs []string
	for _, name := range args.RegularFiles {
		switch {
		case cfg.modules.enabled && strings.HasSuffix(name, ".module.css"):
			moduleSrcs = append(moduleSrcs, name)
		case strings.HasSuffix(strings.ToLower(name), ".css"):
			cssSrcs = append(cssSrcs, name)
		}
	}
	sort.Strings(cssSrcs)
	sort.Strings(moduleSrcs)

	var rules []*rule.Rule
	if len(cssSrcs) > 0 {
		rules = append(rules, newCSSLibrary(args.Rel, cfg, cssSrcs))
	}
	if len(moduleSrcs) > 0 {
		r := rule.NewRule(cssModuleKind, cfg.modules.name)
		r.SetAttr("srcs", moduleSrcs)
		rules = append(rules, r)
	}
	imports := make([]interface{}, len(rules))
	return language.GenerateResult{Gen: rules, Imports: imports}
}

func newCSSLibrary(rel string, cfg *cssConfig, srcs []string) *rule.Rule {
	name := cfg.name
	if name == "" {
		name = path.Base(rel)
		if name == "." || name == "" {
			name = "css"
		}
	}
	r := rule.NewRule(cssLibraryKind, name)
	r.SetAttr("srcs", srcs)
	r.SetAttr("visibility", cfg.visibility)
	return r
}
