package database

import (
	"strconv"
	"testing"
)

func FuzzKeywordNameValidator(f *testing.F) {
	k := KeywordProps{}

	f.Add("My keyword name", true)
	f.Add("My  keyword name", false)
	f.Add("My keyword  name", false)
	f.Add("    My keyword name", false)
	f.Add("My keyword name    ", false)
	f.Add("  My keyword name  ", false)
	f.Add(" My  keyword name  ", false)

	f.Fuzz(func(t *testing.T, name string, valid bool) {
		if ok := k.NameValidator(name); ok != valid {
			t.Errorf(
				"expected %s to be %t but got %t instead",
				strconv.Quote(name),
				valid,
				ok,
			)
		}
	})
}
