package routes

import (
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
