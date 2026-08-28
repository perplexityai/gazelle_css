package css

import (
	"reflect"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func TestImportsIndexesMappedCSSRules(t *testing.T) {
	c := &config.Config{KindMap: map[string]config.MappedKind{
		cssLibraryKind: {
			FromKind: cssLibraryKind,
			KindName: "custom_css_library",
			KindLoad: "//tools:css.bzl",
		},
		cssModuleKind: {
			FromKind: cssModuleKind,
			KindName: "custom_css_module_library",
			KindLoad: "//tools:css.bzl",
		},
	}}
	tests := []struct {
		kind string
		src  string
	}{
		{kind: "custom_css_library", src: "global.css"},
		{kind: "custom_css_module_library", src: "Button.module.css"},
	}
	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			r := rule.NewRule(tt.kind, "css")
			r.SetAttr("srcs", []string{tt.src})
			f := &rule.File{Pkg: "components"}

			got := NewLanguage().Imports(c, r, f)
			want := []resolve.ImportSpec{{Lang: languageName, Imp: "components/" + tt.src}}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("Imports() = %v, want %v", got, want)
			}
		})
	}
}
