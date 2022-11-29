package database

import (
	"fmt"
	"strings"

	"github.com/gookit/validate"
)

const (
	minNameChars = 2
	maxNameChars = 120
)

var notAllowedChars = "|"

// IsAlphaNumeric Check if string contains only alhanumeric
// characters (including underscore).
func IsAlphaNumeric(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') &&
			(r < 'A' || r > 'Z') &&
			(r < '0' || r > '9') &&
			r != '_' && r != ' ' {
			return false
		}
	}
	return true
}

func StartsEndsWithSpace(s string) bool {
	return s[0] == ' ' || s[len(s)-1] == ' '
}

func StartsWithUnderscore(s string) bool {
	return s[0] == '_'
}

// HasMultiSpaceSeparation Check if s has more than 1 space separation
// between words inside it.
func HasMultiSpaceSeparation(s string) bool {
	return strings.Count(s, " ") != len(strings.Fields(s))-1
}

// StartsWithNumber Check if s starts with a number (positive or negative)
func StartsWithNumber(s string) bool {
	i := 0
	n, _ := fmt.Sscanf(s, "%d", &i)
	return n > 0
}

func (k KeywordProps) NameValidator(value string) bool {
	if len(value) < minNameChars || len(value) > maxNameChars {
		return false
	}

	if !IsAlphaNumeric(value) {
		return false
	}

	if StartsWithUnderscore(value) {
		return false
	}

	if StartsEndsWithSpace(value) {
		return false
	}

	if HasMultiSpaceSeparation(value) {
		return false
	}

	if StartsWithNumber(value) {
		return false
	}

	return true
}

func (k KeywordProps) ArgsValidator(value string) bool {
	if StartsEndsWithSpace(value) {
		return false
	}

	return !strings.ContainsAny(value, notAllowedChars)
}

func (k KeywordProps) DocsValidator(value string) bool {
	return !strings.ContainsAny(value, notAllowedChars)
}

func (k KeywordProps) KwTypeValidator(value string) bool {
	return value == Business || value == Technical
}

// Messages Add custom validator messages to keyword props
func (k KeywordProps) Messages() map[string]string {
	return validate.MS{
		"Name.nameValidator": `Bad format of {field}. 
        Expected alphanumeric symbols (incl. underscore). 
        Leading and trailing spaces are not allowed. 
        Multiple word separating spaces are not allowed.`,

		"Args.argsValidator": `Bad format of {field}. 
        '|' symbol is not allowed. 
        Leading and trailing spaces are not allowed.`,

		"Docs.docsValidator": `Bad format of {field}. 
        '|' symbol is not allowed.`,

		"KwType.kwTypeValidator": `Bad format of {field}. 
        Must be one of 'business', 'technical'.`,
	}
}
