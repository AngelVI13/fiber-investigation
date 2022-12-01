package routes

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/AngelVI13/fiber-investigation/pkg/session"
	"github.com/AngelVI13/fiber-investigation/pkg/testutil"
	"github.com/gofiber/fiber/v2"
)

type testHandler func(app *fiber.App, t *testing.T)

func NewTestRouter(t *testing.T) *Router {
	db := testutil.NewTestDb(t)
	return NewRouter(db)
}

func TestRouter(t *testing.T) {
	app := testutil.NewTestFiberApp(t)
	session.CreateSession()

	// closure to provide app and router to testing func
	withArgs := func(h testHandler) func(t *testing.T) {
		return func(t *testing.T) {
			h(app, t)
		}
	}

	t.Log("\n--------Setup done for API testing--------\n")

	t.Run(ExportCsvUrl, withArgs(VerifyExportCsvGet))
	t.Run(ExportStubsUrl, withArgs(VerifyExportStubsGet))
	t.Run(CreateKwdUrl, withArgs(VerifyCreateKeywordPost))
}

func VerifyExportCsvGet(app *fiber.App, t *testing.T) {
	router := NewTestRouter(t)
	app.Get(ExportCsvUrl, Handler(router.HandleExportCsvGet))

	r := httptest.NewRequest(http.MethodGet, ExportCsvUrl, http.NoBody)
	resp, err := app.Test(r, -1)
	if err != nil {
		t.Fatalf("app test request error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
}

func VerifyExportStubsGet(app *fiber.App, t *testing.T) {
	router := NewTestRouter(t)
	app.Get(ExportStubsUrl, Handler(router.HandleExportStubsGet))

	r := httptest.NewRequest(http.MethodGet, ExportStubsUrl, http.NoBody)
	resp, err := app.Test(r, -1)
	if err != nil {
		t.Fatalf("app test request error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("unexpected status code %d", resp.StatusCode)
	}
}

func VerifyCreateKeywordPost(app *fiber.App, t *testing.T) {
	router := NewTestRouter(t)
	app.Get(CreateKwdUrl, Handler(router.HandleCreateKeywordPost))

	initialKeywords, err := database.AllKeywords(router.db)
	if err != nil {
		t.Fatalf("error while getting keywords: %v", err)
	}

	initialKeywordsNum := 6
	if len(initialKeywords) != initialKeywordsNum {
		t.Fatalf(
			"expected %d keywords but got %d",
			initialKeywordsNum,
			len(initialKeywords),
		)
	}

	// TODO: Update this
	var jsonStr = []byte(`{"username": "aditira", "password": "1234"}`)
	r := httptest.NewRequest(http.MethodPost, CreateKwdUrl, bytes.NewBuffer(jsonStr))

	resp, err := app.Test(r, -1)
	if err != nil {
		t.Fatalf("app test request error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Log(resp)
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}

	keywords, err := database.AllKeywords(router.db)
	if err != nil {
		t.Fatalf("error while getting keywords: %v", err)
	}

	if len(keywords) != initialKeywordsNum+1 {
		t.Fatalf(
			"expected %d keywords but got %d",
			initialKeywordsNum+1,
			len(keywords),
		)
	}
}
