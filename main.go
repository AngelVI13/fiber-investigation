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

var UrlMap = map[string]string{
	"BusinessKwdsUrl": "/business",
}

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
		return c.Render("views/index", UpdateFiberMap(UrlMap, fiber.Map{
			"Title":     "Hello, World!",
			"Items":     items,
			"MoreItems": moreItems,
		}), "views/layouts/main")
	})

	app.Get("/:name", indexNameHandler)

	log.Fatal(app.Listen(":3000"))
}

type RandomItem struct {
	Name     string
	Quantity int
}

// UpdateMap update map `n` with values from map `m`
func UpdateFiberMap[T any](m map[string]T, n fiber.Map) fiber.Map {
	for k, v := range m {
		n[k] = v
	}
	return n
}

func indexNameHandler(c *fiber.Ctx) error {
	var items []string

	moreItems := RandomItem{
		Name:     "Watch",
		Quantity: 10,
	}

	// Render index - start with views directory
	return c.Render("views/index", UpdateFiberMap(UrlMap, fiber.Map{
		"Title":     fmt.Sprintf("Hello, %s!", c.Params("name")),
		"Items":     items,
		"MoreItems": moreItems,
	}))
}
