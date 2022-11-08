package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/AngelVI13/fiber-investigation/pkg/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

//go:embed views/*
var viewsfs embed.FS

func main() {
	db_path := "test.db"
	db, err := database.Create(db_path)
	if err != nil {
		log.Fatalf("Couldn't open database: %s - %v", db_path, err)
	}

	engine := html.NewFileSystem(http.FS(viewsfs), ".html")

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	v := routes.NewIndexRoute("/", "views/index", db, []string{"views/layouts/main"})
	app.Get("/", v.Handle)

	log.Fatal(app.Listen(":3000"))
}
