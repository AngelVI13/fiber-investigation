package routes

import (
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/AngelVI13/fiber-investigation/pkg/dbtest"
	"github.com/AngelVI13/fiber-investigation/pkg/session"
)

func TestGenerateCsvFile(t *testing.T) {
	now := time.Now()
	keyword := database.Keyword{
		KeywordProps: database.KeywordProps{
			ValidFrom:      now,
			ValidTo:        nil,
			Name:           "My keyword name",
			Args:           "arg1='a', arg2='b'",
			Docs:           "Very important docstring",
			KwType:         "",
			Implementation: "",
		},
	}

	filename, err := generateCsvFile("test_keywords.csv", []database.Keyword{keyword})
	if err != nil {
		t.Fatalf("error while generating csv file: %v", err)
	}

	if filename == "" {
		t.Fatalf("no filename provided")
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("generated csv file doesn't exist: %v", err)
	}
}

func TestGenerateCsv(t *testing.T) {
	now := time.Now()
	keyword := database.Keyword{
		KeywordProps: database.KeywordProps{
			ValidFrom:      now,
			ValidTo:        nil,
			Name:           "My keyword name",
			Args:           "arg1='a', arg2='b'",
			Docs:           "Very important docstring",
			KwType:         "",
			Implementation: "",
		},
	}

	expCsv := `Name|Args|Docs|Type|Implementation
My keyword name|arg1='a', arg2='b'|Very important docstring||
`

	generatedCsv, err := generateCsv([]database.Keyword{keyword})
	if err != nil {
		t.Fatalf("error while generating csv: %v", err)
	}

	if generatedCsv != expCsv {
		t.Fatalf(
			"mismatch between expected and generated csv:\nexpected:\n%s\nactual:\n%s",
			strconv.Quote(expCsv),
			strconv.Quote(generatedCsv),
		)
	}
}

func TestFiberApp(t *testing.T) {
	app, db := dbtest.NewFiberTest(t)
	router := NewRouter(db)
	session.CreateSession()

	app.Get(ExportCsvUrl, Handler(router.HandleExportCsvGet))

	r := httptest.NewRequest("GET", ExportCsvUrl, nil)
	resp, _ := app.Test(r, -1)
	if resp.StatusCode != 200 {
		t.Errorf("unexpected status code %d", resp.StatusCode)
	}
}
