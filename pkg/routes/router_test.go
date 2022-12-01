package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVI13/fiber-investigation/pkg/dbtest"
	"github.com/AngelVI13/fiber-investigation/pkg/session"
	"github.com/gofiber/fiber/v2"
)

type testHandler func(app *fiber.App, router *Router, t *testing.T)

func TestRouter(t *testing.T) {
	app, db := dbtest.NewTestFiberApp(t)
	// TODO: create a new DB&router for each separate test so that api tests
	// don't influence eachother
	router := NewRouter(db)
	session.CreateSession()

	// closure to provide app and router to testing func
	withArgs := func(h testHandler) func(t *testing.T) {
		return func(t *testing.T) {
			h(app, router, t)
		}
	}

	t.Log("\n--------Setup done for API testing--------\n")

	t.Run(ExportCsvUrl, withArgs(VerifyExportCsvGet))
	t.Run(ExportStubsUrl, withArgs(VerifyExportStubsGet))
}

func VerifyExportCsvGet(app *fiber.App, router *Router, t *testing.T) {
	app.Get(ExportCsvUrl, Handler(router.HandleExportCsvGet))

	r := httptest.NewRequest("GET", ExportCsvUrl, http.NoBody)
	resp, _ := app.Test(r, -1)
	t.Log(resp)

	if resp.StatusCode != 201 {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
}

func VerifyExportStubsGet(app *fiber.App, router *Router, t *testing.T) {
	app.Get(ExportStubsUrl, Handler(router.HandleExportStubsGet))

	r := httptest.NewRequest("GET", ExportStubsUrl, http.NoBody)
	resp, _ := app.Test(r, -1)

	if resp.StatusCode != 200 {
		t.Errorf("unexpected status code %d", resp.StatusCode)
	}
}
