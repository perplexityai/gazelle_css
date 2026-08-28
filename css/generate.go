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
	r := rule.NewRule(cssLibraryKind, name)
	r.SetAttr("srcs", srcs)
	r.SetAttr("visibility", cfg.visibility)
	// Gazelle requires one import-data entry per generated rule. CSS import
	// resolution is intentionally deferred, so the corresponding value is nil.
	return language.GenerateResult{
		Gen:     []*rule.Rule{r},
		Imports: []interface{}{nil},
	}
}
