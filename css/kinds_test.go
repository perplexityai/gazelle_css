package css

import "testing"

func TestCSSModuleKindContract(t *testing.T) {
	info, ok := NewLanguage().Kinds()[cssModuleKind]
	if !ok {
		t.Fatalf("Kinds() does not declare %q", cssModuleKind)
	}
	if !info.NonEmptyAttrs["name"] || !info.NonEmptyAttrs["srcs"] {
		t.Fatalf("%s non-empty attrs = %v, want name and srcs", cssModuleKind, info.NonEmptyAttrs)
	}
	if !info.MergeableAttrs["srcs"] {
		t.Fatalf("%s srcs are not Gazelle-managed", cssModuleKind)
	}

	for _, load := range NewLanguage().Loads() {
		for _, symbol := range load.Symbols {
			if symbol == cssModuleKind {
				return
			}
		}
	}
	t.Fatalf("Loads() does not expose %q", cssModuleKind)
}
