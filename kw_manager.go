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
		Views: engine,
	})

	router := routes.NewRouter(db)
	app.Get(routes.UrlMap["IndexUrl"], router.HandleIndex)
	app.Get(routes.UrlMap["BusinessKwdsUrl"], router.HandleBusinessKeywords)
	app.Get(routes.UrlMap["TechnicalKwdsUrl"], router.HandleTechnicalKeywords)
	app.Get(routes.UrlMap["AllKwdsUrl"], router.HandleAllKeywords)

	app.Get(fmt.Sprintf("%s/:kw_type", routes.UrlMap["CreateKwdUrl"]), router.HandleCreateKeywordGet)
	app.Post(fmt.Sprintf("%s/:kw_type", routes.UrlMap["CreateKwdUrl"]), router.HandleCreateKeywordPost)

	app.Get(fmt.Sprintf("%s/:id", routes.UrlMap["EditKwdUrl"]), router.HandleEditKeywordGet)
	app.Post(fmt.Sprintf("%s/:id", routes.UrlMap["EditKwdUrl"]), router.HandleEditKeywordPost)

	app.Get(fmt.Sprintf("%s/:id", routes.UrlMap["DeleteKwdUrl"]), router.HandleDeleteKeyword)

	app.Get(routes.UrlMap["ExportCsvUrl"], router.HandleExportCsv)
	app.Get(routes.UrlMap["ExportStubsUrl"], router.HandleExportStubsGet)
	app.Post(routes.UrlMap["ExportStubsUrl"], router.HandleExportStubsPost)

	app.Get(routes.UrlMap["ChangelogUrl"], router.HandleChangelog)

	app.Get("/:kwType/version/:id", router.HandleKeywordVersion)

	log.Fatal(app.Listen(":3000"))
}
