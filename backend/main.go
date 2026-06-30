package main

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed web/static
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

	e := newRouter(staticFS)
	log.Fatal(e.Start(":8000"))
}
