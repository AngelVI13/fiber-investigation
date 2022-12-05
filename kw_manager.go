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

	app.Get(routes.IndexUrl, routes.Handler(router.HandleIndex))

	app.Get(routes.BusinessKwdsUrl, routes.Handler(router.HandleBusinessKeywords))
	app.Get(routes.TechnicalKwdsUrl, routes.Handler(router.HandleTechnicalKeywords))
	app.Get(routes.AllKwdsUrl, routes.Handler(router.HandleAllKeywords))

	app.Get(routes.CreateKwdUrlFull, auth.RolesRequires(database.RoleAdmin), routes.Handler(router.HandleCreateKeywordGet))
	app.Post(routes.CreateKwdUrlFull, auth.RolesRequires(database.RoleAdmin), routes.Handler(router.HandleCreateKeywordPost))

	app.Get(routes.EditKwdUrlFull, auth.RolesRequires(database.RoleAdmin), routes.Handler(router.HandleEditKeywordGet))
	app.Post(routes.EditKwdUrlFull, auth.RolesRequires(database.RoleAdmin), routes.Handler(router.HandleEditKeywordPost))

	app.Get(routes.DeleteKwdUrlFull, auth.RolesRequires(database.RoleAdmin), routes.Handler(router.HandleDeleteKeyword))

	app.Get(routes.ImportCsvUrl, auth.RolesRequires(database.RoleAdmin), routes.Handler(router.HandleImportCsvGet))
	app.Post(routes.ImportCsvUrl, auth.RolesRequires(database.RoleAdmin), routes.Handler(router.HandleImportCsvPost))

	app.Get(routes.ExportCsvUrl, routes.Handler(router.HandleExportCsvGet))
	app.Post(routes.ExportCsvUrl, routes.Handler(router.HandleExportCsvPost))

	app.Get(routes.ExportStubsUrl, routes.Handler(router.HandleExportStubsGet))
	app.Post(routes.ExportStubsUrl, routes.Handler(router.HandleExportStubsPost))

	app.Get(routes.ChangelogUrl, routes.Handler(router.HandleChangelog))

	app.Get("/:kwType/version/:id", routes.Handler(router.HandleKeywordVersion))

	app.Get(routes.RegisterUserUrl, auth.RolesRequires(database.RoleAnonimous), routes.Handler(router.HandleRegisterGet))
	app.Post(routes.RegisterUserUrl, auth.RolesRequires(database.RoleAnonimous), routes.Handler(router.HandleRegisterPost))

	app.Get(routes.LoginUrl, auth.RolesRequires(database.RoleAnonimous), routes.Handler(router.HandleLoginGet))
	app.Post(routes.LoginUrl, auth.RolesRequires(database.RoleAnonimous), routes.Handler(router.HandleLoginPost))

	app.Get(routes.LogoutUrl, auth.RolesRequires(database.RoleUser, database.RoleAdmin), routes.Handler(router.HandleLogout))

	app.Get(routes.AdminPanelUrl, auth.RolesRequires(database.RoleAdmin), routes.Handler(router.HandleAdmin))

	app.Get(routes.UserPanelUrl, auth.RolesRequires(database.RoleUser, database.RoleAdmin), routes.Handler(router.HandleUserPanelGet))
	app.Post(routes.UserPanelUrl, auth.RolesRequires(database.RoleUser, database.RoleAdmin), routes.Handler(router.HandleUserPanelPost))

	deleteUserUrl := fmt.Sprintf("%s/:username", routes.DeleteUserUrl)
	editUserUrl := fmt.Sprintf("%s/:username", routes.EditUserUrl)

	app.Get(deleteUserUrl, auth.RolesRequires(database.RoleAdmin), routes.Handler(router.HandleDeleteUser))
	app.Get(editUserUrl, auth.RolesRequires(database.RoleAdmin), routes.Handler(router.HandleEditUserGet))
	app.Post(editUserUrl, auth.RolesRequires(database.RoleAdmin), routes.Handler(router.HandleEditUserPost))

	app.Get(routes.AddUserUrl, auth.RolesRequires(database.RoleAdmin), routes.Handler(router.HandleAddUserGet))
	app.Post(routes.AddUserUrl, auth.RolesRequires(database.RoleAdmin), routes.Handler(router.HandleAddUserPost))

	log.Fatal(app.Listen(":3000"))
}
