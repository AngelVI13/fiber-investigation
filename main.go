package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

//go:embed views/*
var viewsfs embed.FS

func main() {
	engine := html.NewFileSystem(http.FS(viewsfs), ".html")

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	items := []string{"Banana", "Cucumber", "Avocado"}

	moreItems := RandomItem{
		Name:     "Watch",
		Quantity: 10,
	}

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index - start with views directory
		return c.Render("views/index", fiber.Map{
			"Title":     "Hello, World!",
			"Items":     items,
			"MoreItems": moreItems,
		}, "views/layouts/main")
	})

	app.Get("/:name", indexNameHandler)

	log.Fatal(app.Listen(":3000"))
}

type RandomItem struct {
	Name     string
	Quantity int
}

func indexNameHandler(c *fiber.Ctx) error {
	var items []string

	moreItems := RandomItem{
		Name:     "Watch",
		Quantity: 10,
	}

	// Render index - start with views directory
	return c.Render("views/index", fiber.Map{
		"Title":     fmt.Sprintf("Hello, %s!", c.Params("name")),
		"Items":     items,
		"MoreItems": moreItems,
	})
}
