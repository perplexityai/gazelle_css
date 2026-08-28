package css

import "github.com/bazelbuild/bazel-gazelle/rule"

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
