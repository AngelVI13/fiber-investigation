package routes

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/AngelVI13/fiber-investigation/pkg/session"
	"github.com/AngelVI13/fiber-investigation/pkg/testutil"
	"github.com/gofiber/fiber/v2"
)

const (
	FiberCookieName = "fiber-app-flash"
)

type testHandler func(app *fiber.App, t *testing.T)

func NewTestRouter(t *testing.T) *Router {
	db := testutil.NewTestDb(t)
	return NewRouter(db)
}

func MultipartForm(data map[string]string) (
	body *bytes.Buffer,
	formType string,
	err error,
) {
	// New multipart writer.
	body = &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	for k, v := range data {
		fw, err := writer.CreateFormField(k)
		if err != nil {
			return nil, "", fmt.Errorf(
				"failed to create form field %s: %v", k, err,
			)
		}
		_, err = io.Copy(fw, strings.NewReader(v))
		if err != nil {
			return nil, "", fmt.Errorf(
				"failed to add value %s to form field %s: %v", v, k, err,
			)
		}
	}

	return body, writer.FormDataContentType(), nil
}

func CheckCookieStatus(resp *http.Response, status string, t *testing.T) {
	cookies := resp.Cookies()
	if len(cookies) < 1 {
		t.Errorf("unexpected number of cookies: wanted 1 but got %d", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != FiberCookieName {
		t.Errorf("unexpected cookied name: %s", cookie.Name)
	}

	if !strings.Contains(cookie.Value, status) {
		t.Fatalf("expected %s to be flashed on screen but got: %s", status, cookie.Value)
	}
}

func MakePostRequest(
	app *fiber.App,
	url string,
	formData map[string]string,
) (*http.Response, error) {
	var req *http.Request

	if len(formData) == 0 {
		// No form data
		req = httptest.NewRequest(http.MethodPost, url, http.NoBody)
	} else {
		// Build and attach form data to request
		formBody, formType, err := MultipartForm(formData)
		if err != nil {
			return nil, fmt.Errorf("failed to build multipart form: %v", err)
		}

		body := bytes.NewReader(formBody.Bytes())
		req = httptest.NewRequest(http.MethodPost, url, body)
		req.Header.Set("Content-Type", formType)
	}

	resp, err := app.Test(req, -1)
	if err != nil {
		return nil, fmt.Errorf("app test request error: %v", err)
	}

	return resp, nil
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
	t.Run(CreateKwdUrl+" missing name", withArgs(VerifyCreateKeywordMissingNamePost))
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

	data := map[string]string{
		"name": "New keyword",
		"args": "arg1=5, arg2=10",
		"docs": "New keyword docs.",
	}

	url := fmt.Sprintf("%s/%s", CreateKwdUrl, database.Technical)
	resp, err := MakePostRequest(app, url, data)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 302 {
		t.Errorf("unexpected status code %d", resp.StatusCode)
	}

	CheckCookieStatus(resp, "success", t)

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

func VerifyCreateKeywordMissingNamePost(app *fiber.App, t *testing.T) {
	router := NewTestRouter(t)
	app.Post(CreateKwdUrlFull, Handler(router.HandleCreateKeywordPost))

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

	data := map[string]string{
		"name": "",
		"args": "arg1=5, arg2=10",
		"docs": "New keyword docs.",
	}

	url := fmt.Sprintf("%s/%s", CreateKwdUrl, database.Technical)
	resp, err := MakePostRequest(app, url, data)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 302 {
		t.Errorf("unexpected status code %d", resp.StatusCode)
	}

	CheckCookieStatus(resp, "error", t)

	keywords, err := database.AllKeywords(router.db)
	if err != nil {
		t.Fatalf("error while getting keywords: %v", err)
	}

	expKeywordsNum := initialKeywordsNum
	if len(keywords) != expKeywordsNum {
		t.Fatalf(
			"expected %d keywords but got %d",
			expKeywordsNum,
			len(keywords),
		)
	}
}
