// Package css implements a Gazelle extension for CSS source packages.
package css

import (
	"log"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	languageName   = "css"
	cssLibraryKind = "css_library"
	cssModuleKind  = "css_module_library"
)

type fatalFunc func(format string, args ...any)

type cssLang struct {
	fatalf fatalFunc
}

// NewLanguage creates the CSS Gazelle language extension.
func NewLanguage() language.Language { return newLanguage(log.Fatalf) }

func newLanguage(fatalf fatalFunc) *cssLang { return &cssLang{fatalf: fatalf} }

func (l *cssLang) Name() string                                        { return languageName }
func (l *cssLang) Embeds(r *rule.Rule, from label.Label) []label.Label { return nil }
