package css

import (
	"path"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
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
	if err := validateCSSModulePackageSources(moduleSrcs); err != nil {
		rel := args.Rel
		if rel == "" {
			rel = "."
		}
		l.fatalf("%s: %v", rel, err)
		return language.GenerateResult{}
	}

	var rules []*rule.Rule
	var empty []*rule.Rule
	libraryName := cssLibraryName(args.Rel, cfg)
	var retiringLibrary *rule.Rule
	if cfg.modules.enabled && len(moduleSrcs) > 0 && len(cssSrcs) == 0 {
		retiringLibrary = existingRule(args.Config, args.File, cssLibraryKind, libraryName)
	}
	if len(moduleSrcs) > 0 || hasExistingCSSModuleRule(args.Config, args.File) {
		generatesLibraryAndModules := len(moduleSrcs) > 0 && len(cssSrcs) > 0
		if err := validateCSSModuleOwnership(args.Config, args.File, args.OtherGen, cfg.modules.name, libraryName, generatesLibraryAndModules, retiringLibrary); err != nil {
			rel := args.Rel
			if rel == "" {
				rel = "."
			}
			l.fatalf("%s: %v", rel, err)
			return language.GenerateResult{}
		}
	}
	if len(cssSrcs) > 0 {
		rules = append(rules, newCSSLibrary(libraryName, cfg, cssSrcs))
	} else if retiringLibrary != nil {
		empty = append(empty, rule.NewRule(cssLibraryKind, libraryName))
	}
	if cfg.modules.enabled && len(moduleSrcs) > 0 {
		r := rule.NewRule(cssModuleKind, cfg.modules.name)
		r.SetAttr("srcs", moduleSrcs)
		rules = append(rules, r)
	} else if existingRule(args.Config, args.File, cssModuleKind, cfg.modules.name) != nil {
		empty = append(empty, rule.NewRule(cssModuleKind, cfg.modules.name))
	}
	imports := make([]interface{}, len(rules))
	return language.GenerateResult{Gen: rules, Empty: empty, Imports: imports}
}

func cssLibraryName(rel string, cfg *cssConfig) string {
	name := cfg.name
	if name == "" {
		name = path.Base(rel)
		if name == "." || name == "" {
			name = "css"
		}
	}
	return name
}

func newCSSLibrary(name string, cfg *cssConfig, srcs []string) *rule.Rule {
	r := rule.NewRule(cssLibraryKind, name)
	r.SetAttr("srcs", srcs)
	r.SetAttr("visibility", cfg.visibility)
	return r
}

func existingRule(c *config.Config, f *rule.File, kind, name string) *rule.Rule {
	if f == nil {
		return nil
	}
	for _, r := range f.Rules {
		if kindMatches(c, r.Kind(), kind) && r.Name() == name {
			return r
		}
	}
	return nil
}

func kindMatches(c *config.Config, actual, canonical string) bool {
	if actual == canonical {
		return true
	}
	if c == nil {
		return false
	}
	mapped, ok := c.KindMap[canonical]
	return ok && actual == mapped.KindName
}
