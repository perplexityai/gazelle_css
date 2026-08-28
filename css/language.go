// Package css implements a Gazelle extension for CSS source packages.
package css

import (
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	languageName   = "css"
	cssLibraryKind = "css_library"
	cssModuleKind  = "css_module_library"
)

type cssLang struct{}

// NewLanguage creates the CSS Gazelle language extension.
func NewLanguage() language.Language { return &cssLang{} }

func (l *cssLang) Name() string                                        { return languageName }
func (l *cssLang) Embeds(r *rule.Rule, from label.Label) []label.Label { return nil }
