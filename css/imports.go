package css

import (
	"path"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

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
