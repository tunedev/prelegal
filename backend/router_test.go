package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

func testStaticFS() fstest.MapFS {
	return fstest.MapFS{
		"index.html":    {Data: []byte("<html>index</html>")},
		"assets/app.js": {Data: []byte("console.log('app')")},
	}
}

func TestRouter_ServesStaticAsset(t *testing.T) {
	e := newRouter(testStaticFS(), newOpenRouterClient("test-key"))

	req := httptest.NewRequest(http.MethodGet, "/assets/app.js", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "console.log('app')" {
		t.Errorf("expected asset contents, got %q", rec.Body.String())
	}
}

func TestRouter_FallsBackToIndexForClientRoutes(t *testing.T) {
	e := newRouter(testStaticFS(), newOpenRouterClient("test-key"))

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "<html>index</html>" {
		t.Errorf("expected index.html fallback, got %q", rec.Body.String())
	}
}

func TestRouter_HealthRouteNotOverriddenByStatic(t *testing.T) {
	e := newRouter(testStaticFS(), newOpenRouterClient("test-key"))

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	expected := `{"status":"ok"}` + "\n"
	if rec.Body.String() != expected {
		t.Errorf("expected health response, got %q", rec.Body.String())
	}
}
