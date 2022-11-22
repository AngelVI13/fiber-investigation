package routes

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
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
		t.Errorf("error while generating csv file: %v", err)
	}

	if filename == "" {
		t.Errorf("no filename provided")
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("generated csv file doesn't exist: %v", err)
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
		t.Errorf("error while generating csv: %v", err)
	}

	if generatedCsv != expCsv {
		t.Errorf(
			"mismatch between expected and generated csv:\nexpected:\n%s\nactual:\n%s",
			strconv.Quote(expCsv),
			strconv.Quote(generatedCsv),
		)
	}
}
