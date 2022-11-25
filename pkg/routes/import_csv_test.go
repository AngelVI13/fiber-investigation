package routes

import (
	"strconv"
	"testing"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
)

func TestKeywordsFromCsv(t *testing.T) {
	csvText := `Name|Args|Docs|Type|Implementation|
My new keyword|"arg1=""a"", arg2=""b"""|This is my special keyword docs|business||`

	keywords, err := keywordsFromCsv(csvText)
	if err != nil {
		t.Fatal(err)
	}

	if len(keywords) != 1 {
		t.Fatalf("expected 1 keyword but got %d", len(keywords))
	}

	keyword := keywords[0]

	expName := "My new keyword"
	if keyword.Name != expName {
		t.Errorf(
			"expected %s keyword name but got %s",
			strconv.Quote(expName),
			strconv.Quote(keyword.Name),
		)
	}

	expArgs := "arg1=\"a\", arg2=\"b\""
	if keyword.Args != expArgs {
		t.Errorf(
			"expected %s keyword args but got %s",
			strconv.Quote(expArgs),
			strconv.Quote(keyword.Args),
		)
	}

	expDocs := "This is my special keyword docs"
	if keyword.Docs != expDocs {
		t.Errorf(
			"expected %s keyword docs but got %s",
			strconv.Quote(expDocs),
			strconv.Quote(keyword.Docs),
		)
	}

	expKwType := database.Business
	if keyword.KwType != expKwType {
		t.Errorf(
			"expected %s keyword type but got %s",
			strconv.Quote(expKwType),
			strconv.Quote(keyword.KwType),
		)
	}

	expImplementation := ""
	if keyword.Implementation != expImplementation {
		t.Errorf(
			"expected %s keyword implementation but got %s",
			strconv.Quote(expImplementation),
			strconv.Quote(keyword.Implementation),
		)
	}
}
