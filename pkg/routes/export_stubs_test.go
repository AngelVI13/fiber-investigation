package routes

import "testing"


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

