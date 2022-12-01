package routes

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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
	app.Post(CreateKwdUrlFull, Handler(router.HandleCreateKeywordPost))

	// TODO: should we check for num of keywords before CreateKw request
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

	var jsonStr = []byte(`{"name": "New+keyword", "args": "arg1=5,+arg2=10", "docs": "New+keywords+doc"}`)
	url := fmt.Sprintf("%s/%s", CreateKwdUrl, database.Technical)
	r := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))

	resp, err := app.Test(r, -1)
	if err != nil {
		t.Fatalf("app test request error: %v", err)
	}

	if resp.StatusCode != 302 {
		t.Errorf("unexpected status code %d", resp.StatusCode)
	}

	// TODO: Check len of cookies and cookie name matches
	cookies := resp.Cookies()
	cookie := cookies[0]

	if !strings.Contains(cookie.Value, "success") {
		t.Fatalf("expected success to be flashed on screen but got: %s", cookie.Value)
	}

	keywords, err := database.AllKeywords(router.db)
	if err != nil {
		t.Fatalf("error while getting keywords: %v", err)
	}

	expKeywordsNum := initialKeywordsNum + 1
	if len(keywords) != expKeywordsNum {
		t.Fatalf(
			"expected %d keywords but got %d",
			expKeywordsNum,
			len(keywords),
		)
	}

	foundIdx := -1
	for i, kw := range keywords {
		if kw.Name == "New keyword" {
			foundIdx = i
		}
	}

	if foundIdx == -1 {
		lstKw := keywords[len(keywords)-1]
		t.Log(lstKw)
		t.Log(lstKw.Name, lstKw.Args, lstKw.Docs)
		t.Errorf("did not find newly created keyword: %v", keywords)
	}
}
