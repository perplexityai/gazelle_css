package css

import "testing"

func TestCSSModuleKindContract(t *testing.T) {
	kinds := NewLanguage().Kinds()
	libraryInfo := kinds[cssLibraryKind]
	if len(libraryInfo.NonEmptyAttrs) != 1 || !libraryInfo.NonEmptyAttrs["srcs"] {
		t.Fatalf("%s non-empty attrs = %v, want only srcs", cssLibraryKind, libraryInfo.NonEmptyAttrs)
	}

	info, ok := kinds[cssModuleKind]
	if !ok {
		t.Fatalf("Kinds() does not declare %q", cssModuleKind)
	}
	if len(info.NonEmptyAttrs) != 1 || !info.NonEmptyAttrs["srcs"] {
		t.Fatalf("%s non-empty attrs = %v, want only srcs", cssModuleKind, info.NonEmptyAttrs)
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
