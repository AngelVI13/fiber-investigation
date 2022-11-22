package routes

import (
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
)

func FuzzRfStubGeneratorName(f *testing.F) {
	rfStubGenerator := RfStubGenerator{}

	f.Add("My keyword name", "My keyword name")
	f.Add("My  keyword name", "My keyword name")
	f.Add("My keyword  name", "My keyword name")
	f.Add("    My keyword name", "My keyword name")
	f.Add("My keyword name    ", "My keyword name")
	f.Add("  My keyword name  ", "My keyword name")
	f.Add(" My  keyword name  ", "My keyword name")

	f.Fuzz(func(t *testing.T, rawName string, expName string) {
		cleanName := rfStubGenerator.Name(rawName)

		if cleanName != expName {
			t.Errorf("wrong rf name: expected '%s' but got '%s'", expName, cleanName)
		}
	})
}

func TestRfStubGeneratorHeader(t *testing.T) {
	rfStubGenerator := RfStubGenerator{}

	header := rfStubGenerator.Header()

	if !strings.Contains(header, "*** Keywords ***") {
		t.Errorf("missing keyword section in rf header: %s", header)
	}
}

func TestRfStubGeneratorDocs(t *testing.T) {
	rfStubGenerator := RfStubGenerator{}

	rawDocs := `
Switches the device on and off in order to clear last (previous) program.

If system error occurs and the device is restarted, it attempts to boot to the same 
    program that was running before the error. 
This method is used for clearing previous program in order to boot to standby.


    `
	expDocs := `Switches the device on and off in order to clear last (previous) program.
    ...    
    ...    If system error occurs and the device is restarted, it attempts to boot to the same 
    ...        program that was running before the error. 
    ...    This method is used for clearing previous program in order to boot to standby.`

	cleanDocs := rfStubGenerator.Docs(rawDocs)

	compareDocs(expDocs, cleanDocs, t)
}

func TestPyStubGeneratorDocs(t *testing.T) {
	pyStubGenerator := PyStubGenerator{}

	rawDocs := `
Switches the device on and off in order to clear last (previous) program.

If system error occurs and the device is restarted, it attempts to boot to the same 
    program that was running before the error. 
This method is used for clearing previous program in order to boot to standby.


    `
	expDocs := `Switches the device on and off in order to clear last (previous) program.
    
    If system error occurs and the device is restarted, it attempts to boot to the same 
        program that was running before the error. 
    This method is used for clearing previous program in order to boot to standby.`

	cleanDocs := pyStubGenerator.Docs(rawDocs)

	compareDocs(expDocs, cleanDocs, t)
}

func FuzzPyStubGeneratorRawName(f *testing.F) {
	pyStubGenerator := PyStubGenerator{}

	f.Add("My keyword name", "My keyword name")
	f.Add("My  keyword name", "My keyword name")
	f.Add("My keyword  name", "My keyword name")
	f.Add("    My keyword name", "My keyword name")
	f.Add("My keyword name    ", "My keyword name")
	f.Add("  My keyword name  ", "My keyword name")
	f.Add(" My  keyword name  ", "My keyword name")

	f.Fuzz(func(t *testing.T, rawName string, expName string) {
		cleanName := pyStubGenerator.RawName(rawName)

		if cleanName != expName {
			t.Errorf("wrong py raw name: expected '%s' but got '%s'", expName, cleanName)
		}
	})
}

func FuzzPyStubGeneratorName(f *testing.F) {
	pyStubGenerator := PyStubGenerator{}

	f.Add("My keyword name", "my_keyword_name")
	f.Add("My  keyword name", "my_keyword_name")
	f.Add("My keyword  name", "my_keyword_name")
	f.Add("    My keyword name", "my_keyword_name")
	f.Add("My keyword name    ", "my_keyword_name")
	f.Add("  My keyword name  ", "my_keyword_name")
	f.Add(" My  keyword name  ", "my_keyword_name")

	f.Fuzz(func(t *testing.T, rawName string, expName string) {
		cleanName := pyStubGenerator.Name(rawName)

		if cleanName != expName {
			t.Errorf("wrong py method name: expected '%s' but got '%s'", expName, cleanName)
		}
	})
}

func TestPyStubGeneratorHeader(t *testing.T) {
	pyStubGenerator := PyStubGenerator{}

	header := pyStubGenerator.Header()

	if !strings.Contains(header, "from robot.api.deco import keyword") {
		t.Errorf("missing keyword import in py header: %s", header)
	}
}

func TestGenerateStubsFileRf(t *testing.T) {
	keyword, _, _ := exampleKeywordAndStubs()
	stubGenerator := &RfStubGenerator{}

	checkGenerateStubsFile(stubGenerator, []database.Keyword{keyword}, t)
}

func TestGenerateStubsFilePy(t *testing.T) {
	keyword, _, _ := exampleKeywordAndStubs()
	stubGenerator := &PyStubGenerator{}

	checkGenerateStubsFile(stubGenerator, []database.Keyword{keyword}, t)
}

func TestGenerateStubsRf(t *testing.T) {
	keyword, _, rfStub := exampleKeywordAndStubs()
	stubGenerator := &RfStubGenerator{}

	checkGenerateStubs(stubGenerator, keyword, rfStub, t)
}

func TestGenerateStubsPy(t *testing.T) {
	keyword, pyStub, _ := exampleKeywordAndStubs()
	stubGenerator := &PyStubGenerator{}

	checkGenerateStubs(stubGenerator, keyword, pyStub, t)
}

func exampleKeywordAndStubs() (keyword database.Keyword, pyStub string, rfStub string) {
	now := time.Now()
	keyword = database.Keyword{
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

	pyStub = `from robot.api.deco import keyword

@keyword("My keyword name")
def my_keyword_name(arg1='a', arg2='b'):
    """
    Very important docstring
    """
    pass

    `

	rfStub = `*** Keywords ***

My keyword name
    [Documentation]  Very important docstring
    [Arguments]      arg1='a', arg2='b'
    Log 	         NOP

    `
	return keyword, pyStub, rfStub
}

func checkGenerateStubs(
	stubGenerator StubGenerator,
	keyword database.Keyword,
	expStub string,
	t *testing.T,
) {
	generatedStub, err := generateStubs(stubGenerator, []database.Keyword{keyword})
	if err != nil {
		t.Errorf("error while generating stubs: %v", err)
	}

	if generatedStub != expStub {
		t.Errorf(
			"mismatch between expected and generated stub:\nexpected:\n%s\nactual:\n%s",
			strconv.Quote(expStub),
			strconv.Quote(generatedStub),
		)
	}
}

func checkGenerateStubsFile(
	stubGenerator StubGenerator,
	keywords []database.Keyword,
	t *testing.T,
) {
	filename, err := generateStubsFile(stubGenerator, keywords)
	if err != nil {
		t.Errorf("error while generating stubs file: %v", err)
	}

	if filename == "" {
		t.Errorf("no filename provided")
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("generated stubs file doesn't exist: %v", err)
	}
}

func compareDocs(exp, act string, t *testing.T) {
	expParts := strings.SplitAfter(exp, "\n")
	cleanParts := strings.SplitAfter(act, "\n")

	if len(expParts) != len(cleanParts) {
		t.Errorf(
			"wrong number of py docs lines: expected %d but got %d",
			len(expParts),
			len(cleanParts),
		)
	}

	for i := range cleanParts {
		expLine := expParts[i]
		cleanLine := cleanParts[i]

		if expLine != cleanLine {
			t.Errorf("wrong format of py docs (line %d): \nexpected: \n%s\ngot:\n%s",
				i,
				strconv.Quote(expLine),
				strconv.Quote(cleanLine),
			)
		}
	}
}
