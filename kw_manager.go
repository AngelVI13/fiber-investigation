package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/AngelVI13/fiber-investigation/pkg/auth"
	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/AngelVI13/fiber-investigation/pkg/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

//go:embed views/*
var viewsfs embed.FS

// Handler Wrapper to convert handler args to expected args by fiber and
// add url map to context.
func Handler(f func(c *routes.Ctx) error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := &routes.Ctx{
			Ctx: c,
		}
		return f(ctx.WithUrls().WithSession())
	}
}

func main() {
	db_path := "test.db"
	db, err := database.Create(db_path)
	if err != nil {
		log.Fatalf("Couldn't open database: %s - %v", db_path, err)
	}

	auth.CreateSession()

	engine := html.NewFileSystem(http.FS(viewsfs), ".html")

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: routes.MainLayoutView,
	})

	app.Static("/css", "./views/static/css")

	router := routes.NewRouter(db)

	app.Get(routes.IndexUrl, Handler(router.HandleIndex))

	app.Get(routes.BusinessKwdsUrl, Handler(router.HandleBusinessKeywords))
	app.Get(routes.TechnicalKwdsUrl, Handler(router.HandleTechnicalKeywords))
	app.Get(routes.AllKwdsUrl, Handler(router.HandleAllKeywords))

	app.Get(fmt.Sprintf("%s/:kw_type", routes.CreateKwdUrl), auth.RolesRequires(database.RoleAdmin), Handler(router.HandleCreateKeywordGet))
	app.Post(fmt.Sprintf("%s/:kw_type", routes.CreateKwdUrl), auth.RolesRequires(database.RoleAdmin), Handler(router.HandleCreateKeywordPost))

	app.Get(fmt.Sprintf("%s/:id/:kw_type", routes.EditKwdUrl), auth.RolesRequires(database.RoleAdmin), Handler(router.HandleEditKeywordGet))
	app.Post(fmt.Sprintf("%s/:id/:kw_type", routes.EditKwdUrl), auth.RolesRequires(database.RoleAdmin), Handler(router.HandleEditKeywordPost))

	app.Get(fmt.Sprintf("%s/:id/:kw_type", routes.DeleteKwdUrl), auth.RolesRequires(database.RoleAdmin), Handler(router.HandleDeleteKeyword))

	app.Get(routes.ImportCsvUrl, auth.RolesRequires(database.RoleAdmin), Handler(router.HandleImportCsvGet))
	app.Post(routes.ImportCsvUrl, auth.RolesRequires(database.RoleAdmin), Handler(router.HandleImportCsvPost))

	app.Get(routes.ExportCsvUrl, Handler(router.HandleExportCsvGet))
	app.Post(routes.ExportCsvUrl, Handler(router.HandleExportCsvPost))

	app.Get(routes.ExportStubsUrl, Handler(router.HandleExportStubsGet))
	app.Post(routes.ExportStubsUrl, Handler(router.HandleExportStubsPost))

	app.Get(routes.ChangelogUrl, Handler(router.HandleChangelog))

	app.Get("/:kwType/version/:id", Handler(router.HandleKeywordVersion))

	app.Get(routes.RegisterUserUrl, auth.RolesRequires(database.RoleAnonimous), Handler(router.HandleRegisterGet))
	app.Post(routes.RegisterUserUrl, auth.RolesRequires(database.RoleAnonimous), Handler(router.HandleRegisterPost))

	app.Get(routes.LoginUrl, auth.RolesRequires(database.RoleAnonimous), Handler(router.HandleLoginGet))
	app.Post(routes.LoginUrl, auth.RolesRequires(database.RoleAnonimous), Handler(router.HandleLoginPost))

	app.Get(routes.LogoutUrl, auth.RolesRequires(database.RoleUser, database.RoleAdmin), Handler(router.HandleLogout))

	app.Get(routes.AdminPanelUrl, auth.RolesRequires(database.RoleAdmin), Handler(router.HandleAdmin))

	app.Get(routes.UserPanelUrl, auth.RolesRequires(database.RoleUser, database.RoleAdmin), Handler(router.HandleUserPanelGet))
	app.Post(routes.UserPanelUrl, auth.RolesRequires(database.RoleUser, database.RoleAdmin), Handler(router.HandleUserPanelPost))

	deleteUserUrl := fmt.Sprintf("%s/:username", routes.DeleteUserUrl)
	editUserUrl := fmt.Sprintf("%s/:username", routes.EditUserUrl)

	app.Get(deleteUserUrl, auth.RolesRequires(database.RoleAdmin), Handler(router.HandleDeleteUser))
	app.Get(editUserUrl, auth.RolesRequires(database.RoleAdmin), Handler(router.HandleEditUserGet))
	app.Post(editUserUrl, auth.RolesRequires(database.RoleAdmin), Handler(router.HandleEditUserPost))

	app.Get(routes.AddUserUrl, auth.RolesRequires(database.RoleAdmin), Handler(router.HandleAddUserGet))
	app.Post(routes.AddUserUrl, auth.RolesRequires(database.RoleAdmin), Handler(router.HandleAddUserPost))

	log.Fatal(app.Listen(":3000"))
}
