package database

import (
	"strconv"
	"testing"
)

func FuzzKeywordNameValidator(f *testing.F) {
	k := KeywordProps{}

	f.Add("My keyword name", true)
	// Extra whitespace separation
	f.Add("My  keyword name", false)
	f.Add("My keyword  name", false)
	// Extra leading/trailing whitespace
	f.Add("    My keyword name", false)
	f.Add("My keyword name    ", false)
	f.Add("  My keyword name  ", false)
	f.Add(" My  keyword name  ", false)
	// Non-alphanumeric symbols in name
	f.Add("My keyword name|", false)
	f.Add("My keyword name?", false)
	f.Add("My keyword name%", false)
	f.Add("My keyword name$", false)
	// Starts with underscore
	f.Add("_My keyword name", false)
	// Max/Min name length
	f.Add("m", false)
	f.Add(`My keyword name is very very 
    long long long long long long long long 
    long long long long long long long long`, false)

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

func FuzzKeywordArgsValidator(f *testing.F) {
	k := KeywordProps{}

	f.Add("My keyword args", true)
	// Extra whitespace separation
	f.Add("My  keyword args", true)
	f.Add("My keyword  args", true)
	// Extra leading/trailing whitespace
	f.Add("    My keyword args", false)
	f.Add("My keyword args    ", false)
	f.Add("  My keyword args  ", false)
	f.Add(" My  keyword args  ", false)
	// Non-alphanumeric symbols in args
	// pipe char is not allowed in args field
	f.Add("My keyword args|", false)
	f.Add("My keyword args?", true)
	f.Add("My keyword args%", true)
	f.Add("My keyword args$", true)
	// Starts with underscore
	f.Add("_My keyword args", true)
	// Max/Min args length
	f.Add("m", true)
	f.Add(`My keyword args is very very 
    long long long long long long long long 
    long long long long long long long long`, true)

	f.Fuzz(func(t *testing.T, args string, valid bool) {
		if ok := k.ArgsValidator(args); ok != valid {
			t.Errorf(
				"expected %s to be %t but got %t instead",
				strconv.Quote(args),
				valid,
				ok,
			)
		}
	})
}

func FuzzKeywordDocsValidator(f *testing.F) {
	k := KeywordProps{}

	f.Add("My keyword docs", true)
	// Extra whitespace separation
	f.Add("My  keyword docs", true)
	f.Add("My keyword  docs", true)
	// Extra leading/trailing whitespace
	f.Add("    My keyword docs", true)
	f.Add("My keyword docs    ", true)
	f.Add("  My keyword docs  ", true)
	f.Add(" My  keyword docs  ", true)
	// Non-alphanumeric symbols in docs
	// pipe char is not allowed in docs field
	f.Add("My keyword docs|", false)
	f.Add("My keyword docs?", true)
	f.Add("My keyword docs%", true)
	f.Add("My keyword docs$", true)
	// Starts with underscore
	f.Add("_My keyword docs", true)
	// Max/Min docs length
	f.Add("m", true)
	f.Add(`My keyword docs is very very 
    long long long long long long long long 
    long long long long long long long long`, true)

	f.Fuzz(func(t *testing.T, docs string, valid bool) {
		if ok := k.DocsValidator(docs); ok != valid {
			t.Errorf(
				"expected %s to be %t but got %t instead",
				strconv.Quote(docs),
				valid,
				ok,
			)
		}
	})
}

func FuzzKeywordTypeValidator(f *testing.F) {
	k := KeywordProps{}

	f.Add("My keyword kwType", false)
	// Extra whitespace separation
	f.Add("My  keyword kwType", false)
	f.Add("My keyword  kwType", false)
	// Extra leading/trailing whitespace
	f.Add("    My keyword kwType", false)
	f.Add("My keyword kwType    ", false)
	f.Add("  My keyword kwType  ", false)
	f.Add(" My  keyword kwType  ", false)
	// Non-alphanumeric symbols in kwType
	// pipe char is not allowed in kwType field
	f.Add("My keyword kwType|", false)
	f.Add("My keyword kwType?", false)
	f.Add("My keyword kwType%", false)
	f.Add("My keyword kwType$", false)
	// Starts with underscore
	f.Add("_My keyword kwType", false)
	// Max/Min kwType length
	f.Add("m", false)
	f.Add(`My keyword kwType is very very 
    long long long long long long long long 
    long long long long long long long long`, false)
	// only business & technical kwTypes are allowed
	f.Add(Business, true)
	f.Add(Technical, true)

	f.Fuzz(func(t *testing.T, kwType string, valid bool) {
		if ok := k.KwTypeValidator(kwType); ok != valid {
			t.Errorf(
				"expected %s to be %t but got %t instead",
				strconv.Quote(kwType),
				valid,
				ok,
			)
		}
	})
}
