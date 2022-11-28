package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/AngelVI13/fiber-investigation/pkg/routes"
	"github.com/AngelVI13/fiber-investigation/pkg/session"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

//go:embed views/*
var viewsfs embed.FS

// Handler Wrapper to convert handler args to expected args by fiber
func Handler(f func(c *routes.Ctx) error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := &routes.Ctx{
			Ctx: c,
		}
		return f(ctx)
	}
}

func main() {
	db_path := "test.db"
	db, err := database.Create(db_path)
	if err != nil {
		log.Fatalf("Couldn't open database: %s - %v", db_path, err)
	}

	session.CreateSession()

	engine := html.NewFileSystem(http.FS(viewsfs), ".html")

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "views/layouts/main",
	})

	app.Static("/css", "./views/static/css")

	router := routes.NewRouter(db)
	app.Get(routes.IndexUrl, Handler(router.HandleIndex))

	app.Get(routes.BusinessKwdsUrl, Handler(router.HandleBusinessKeywords))
	app.Get(routes.TechnicalKwdsUrl, Handler(router.HandleTechnicalKeywords))
	app.Get(routes.AllKwdsUrl, Handler(router.HandleAllKeywords))

	app.Get(fmt.Sprintf("%s/:kw_type", routes.CreateKwdUrl), Handler(router.HandleCreateKeywordGet))
	app.Post(fmt.Sprintf("%s/:kw_type", routes.CreateKwdUrl), Handler(router.HandleCreateKeywordPost))

	app.Get(fmt.Sprintf("%s/:id", routes.EditKwdUrl), Handler(router.HandleEditKeywordGet))
	app.Post(fmt.Sprintf("%s/:id", routes.EditKwdUrl), Handler(router.HandleEditKeywordPost))

	app.Get(fmt.Sprintf("%s/:id", routes.DeleteKwdUrl), Handler(router.HandleDeleteKeyword))

	app.Get(routes.ImportCsvUrl, Handler(router.HandleImportCsvGet))
	app.Post(routes.ImportCsvUrl, Handler(router.HandleImportCsvPost))

	app.Get(routes.ExportCsvUrl, Handler(router.HandleExportCsvGet))
	app.Post(routes.ExportCsvUrl, Handler(router.HandleExportCsvPost))

	app.Get(routes.ExportStubsUrl, Handler(router.HandleExportStubsGet))
	app.Post(routes.ExportStubsUrl, Handler(router.HandleExportStubsPost))

	app.Get(routes.ChangelogUrl, Handler(router.HandleChangelog))

	app.Get("/:kwType/version/:id", Handler(router.HandleKeywordVersion))

	app.Get(routes.RegisterUserUrl, Handler(router.HandleRegisterGet))
	app.Post(routes.RegisterUserUrl, Handler(router.HandleRegisterPost))

    app.Get(routes.LoginUrl, Handler(router.HandleLoginGet))
    app.Post(routes.LoginUrl, Handler(router.HandleLoginPost))

	app.Get(routes.LogoutUrl, Handler(router.HandleLoginGet))


	log.Fatal(app.Listen(":3000"))
}
