package routes

import (
	"strconv"
	"strings"
	"testing"
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

	expParts := strings.SplitAfter(expDocs, "\n")
	cleanParts := strings.SplitAfter(cleanDocs, "\n")

	if len(expParts) != len(cleanParts) {
		t.Errorf(
			"wrong number of rf docs lines: expected %d but got %d",
			len(expParts),
			len(cleanParts),
		)
	}

	for i := range cleanParts {
		expLine := expParts[i]
		cleanLine := cleanParts[i]

		if expLine != cleanLine {
			t.Errorf("wrong format of rf docs (line %d): \nexpected: \n%s\ngot:\n%s",
				i,
				strconv.Quote(expLine),
				strconv.Quote(cleanLine),
			)
		}
	}
}
