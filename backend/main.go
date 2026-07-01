package main

import (
	"embed"
	"io/fs"
	"log"
	"os"
)

// The "all:" prefix is required because go:embed otherwise excludes any
// file or directory whose name starts with "_" or ".", and SvelteKit's
// build output puts every asset under "_app/".
//
//go:embed all:web/static
var embeddedStatic embed.FS

const dbPath = "data/app.db"

func main() {
	staticFS, err := fs.Sub(embeddedStatic, "web/static")
	if err != nil {
		log.Fatalf("loading embedded static assets: %v", err)
	}

	db, err := initDB(dbPath)
	if err != nil {
		log.Fatalf("initializing database: %v", err)
	}
	defer db.Close()

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Println("WARNING: OPENROUTER_API_KEY is not set — the AI chat will fail on every request. " +
			"If you're running this outside Docker, export it or run with `--env-file .env` / `set -a; source .env; set +a`.")
	}
	chatClient := newOpenRouterClient(apiKey)

	e := newRouter(staticFS, chatClient)
	log.Fatal(e.Start(":8000"))
}
