package main

import (
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// newRouter wires up the API routes and serves staticFS as a single-page
// app, falling back to index.html for any path that isn't a real asset.
func newRouter(staticFS fs.FS, chatClient *OpenRouterClient) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/api/health", healthHandler)
	e.POST("/api/chat", newChatHandler(chatClient))

	fileServer := http.FileServer(http.FS(staticFS))
	e.GET("/*", echo.WrapHandler(spaFallback(staticFS, fileServer)))

	return e
}

func spaFallback(staticFS fs.FS, fileServer http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path != "/" {
			if _, err := fs.Stat(staticFS, path[1:]); err != nil {
				r.URL.Path = "/"
			}
		}
		fileServer.ServeHTTP(w, r)
	}
}
