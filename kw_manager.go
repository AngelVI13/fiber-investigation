package main

import (
	"embed"
	"fmt"
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
		Views:       engine,
		ViewsLayout: "views/layouts/main",
	})

	app.Static("/css", "./views/static/css")

	router := routes.NewRouter(db)
	app.Get(routes.IndexUrl, func(c *fiber.Ctx) error {
		return router.HandleIndex(&routes.Ctx{c})
	})
	app.Get(routes.BusinessKwdsUrl, func(c *fiber.Ctx) error {
		return router.HandleBusinessKeywords(&routes.Ctx{c})
	})
	app.Get(routes.TechnicalKwdsUrl, func(c *fiber.Ctx) error {
		return router.HandleTechnicalKeywords(&routes.Ctx{c})
	})
	app.Get(routes.AllKwdsUrl, func(c *fiber.Ctx) error {
		return router.HandleAllKeywords(&routes.Ctx{c})
	})

	app.Get(fmt.Sprintf("%s/:kw_type", routes.CreateKwdUrl), func(c *fiber.Ctx) error {
		return router.HandleCreateKeywordGet(&routes.Ctx{c})
	})
	app.Post(fmt.Sprintf("%s/:kw_type", routes.CreateKwdUrl), func(c *fiber.Ctx) error {
		return router.HandleCreateKeywordPost(&routes.Ctx{c})
	})

	app.Get(fmt.Sprintf("%s/:id", routes.EditKwdUrl), func(c *fiber.Ctx) error {
		return router.HandleEditKeywordGet(&routes.Ctx{c})
	})
	app.Post(fmt.Sprintf("%s/:id", routes.EditKwdUrl), func(c *fiber.Ctx) error {
		return router.HandleEditKeywordPost(&routes.Ctx{c})
	})

	app.Get(fmt.Sprintf("%s/:id", routes.DeleteKwdUrl), func(c *fiber.Ctx) error {
		return router.HandleDeleteKeyword(&routes.Ctx{c})
	})

	app.Get(routes.ImportCsvUrl, func(c *fiber.Ctx) error {
		return router.HandleImportCsvGet(&routes.Ctx{c})
	})
	app.Post(routes.ImportCsvUrl, func(c *fiber.Ctx) error {
		return router.HandleImportCsvPost(&routes.Ctx{c})
	})

	app.Get(routes.ExportCsvUrl, func(c *fiber.Ctx) error {
		return router.HandleExportCsvGet(&routes.Ctx{c})
	})
	app.Post(routes.ExportCsvUrl, func(c *fiber.Ctx) error {
		return router.HandleExportCsvPost(&routes.Ctx{c})
	})

	app.Get(routes.ExportStubsUrl, func(c *fiber.Ctx) error {
		return router.HandleExportStubsGet(&routes.Ctx{c})
	})
	app.Post(routes.ExportStubsUrl, func(c *fiber.Ctx) error {
		return router.HandleExportStubsPost(&routes.Ctx{c})
	})

	app.Get(routes.ChangelogUrl, func(c *fiber.Ctx) error {
		return router.HandleChangelog(&routes.Ctx{c})
	})

	app.Get("/:kwType/version/:id", func(c *fiber.Ctx) error {
		return router.HandleKeywordVersion(&routes.Ctx{c})
	})

	log.Fatal(app.Listen(":3000"))
}
