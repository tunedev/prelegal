package main

import (
	"io/fs"
	"testing"
)

// This guards against accidentally dropping the "all:" prefix from the
// go:embed directive in main.go: go:embed excludes any file or directory
// whose name starts with "_" or ".", but SvelteKit's build output puts
// every real asset under "_app/", so without "all:" the entire frontend
// silently fails to embed (the SPA fallback then masks this by serving
// index.html for every asset request, which still returns HTTP 200).
func TestEmbeddedStatic_IncludesUnderscorePrefixedDirectories(t *testing.T) {
	staticFS, err := fs.Sub(embeddedStatic, "web/static")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := fs.Stat(staticFS, "_app/placeholder.txt"); err != nil {
		t.Fatalf("expected _app/placeholder.txt to be embedded, got error: %v", err)
	}
}
