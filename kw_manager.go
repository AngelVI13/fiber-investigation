package main

import (
	"embed"
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
		ViewsLayout: routes.MainLayoutView,
	})

	app.Static("/css", "./views/static/css")

	router := routes.NewRouter(db)

	app.Get(routes.IndexUrl, routes.Handler(router.HandleIndex))

	app.Get(routes.BusinessKwdsUrl, routes.Handler(router.HandleBusinessKeywords))
	app.Get(routes.TechnicalKwdsUrl, routes.Handler(router.HandleTechnicalKeywords))
	app.Get(routes.AllKwdsUrl, routes.Handler(router.HandleAllKeywords))

	app.Get(routes.CreateKwdUrlFull, routes.Handler(router.HandleCreateKeywordGet))
	app.Post(routes.CreateKwdUrlFull, routes.Handler(router.HandleCreateKeywordPost))

	app.Get(routes.EditKwdUrlFull, routes.Handler(router.HandleEditKeywordGet))
	app.Post(routes.EditKwdUrlFull, routes.Handler(router.HandleEditKeywordPost))

	app.Get(routes.DeleteKwdUrlFull, routes.Handler(router.HandleDeleteKeyword))

	app.Get(routes.ImportCsvUrl, routes.Handler(router.HandleImportCsvGet))
	app.Post(routes.ImportCsvUrl, routes.Handler(router.HandleImportCsvPost))

	app.Get(routes.ExportCsvUrl, routes.Handler(router.HandleExportCsvGet))
	app.Post(routes.ExportCsvUrl, routes.Handler(router.HandleExportCsvPost))

	app.Get(routes.ExportStubsUrl, routes.Handler(router.HandleExportStubsGet))
	app.Post(routes.ExportStubsUrl, routes.Handler(router.HandleExportStubsPost))

	app.Get(routes.ChangelogUrl, routes.Handler(router.HandleChangelog))

	app.Get("/:kwType/version/:id", routes.Handler(router.HandleKeywordVersion))

	app.Get(routes.RegisterUserUrl, routes.Handler(router.HandleRegisterGet))
	app.Post(routes.RegisterUserUrl, routes.Handler(router.HandleRegisterPost))

	app.Get(routes.LoginUrl, routes.Handler(router.HandleLoginGet))
	app.Post(routes.LoginUrl, routes.Handler(router.HandleLoginPost))

	app.Get(routes.LogoutUrl, routes.Handler(router.HandleLogout))

	log.Fatal(app.Listen(":3000"))
}
