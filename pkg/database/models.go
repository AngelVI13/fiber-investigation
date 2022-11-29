package database

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
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

type KeywordProps struct {
	ValidFrom      time.Time  `gorm:"autoCreateTime;not null" csv:"-" validate:"-"`
	ValidTo        *time.Time `csv:"-" validate:"-"`
	Name           string     `gorm:"not null" csv:"Name" validate:"nameValidator"`
	Args           string     `gorm:"not null" csv:"Args" validate:"argsValidator"`
	Docs           string     `gorm:"not null" csv:"Docs" validate:"docsValidator"`
	KwType         string     `gorm:"not null" csv:"Type" validate:"kwTypeValidator"`
	Implementation string     `csv:"Implementation" validate:"-"`
}

func (k KeywordProps) String() string {
	return fmt.Sprintf("%s", k.Name)
}

func (k KeywordProps) NameValidator(value string) bool {
	if len(value) < minNameChars || len(value) > maxNameChars {
		return false
	}

	if !IsAlphaNumeric(value) {
		return false
	}

	// check if name starts with underscore or empty space
	if value[0] == '_' {
		return false
	}

	// check if name starts/ends with space
	if value[0] == ' ' || value[len(value)-1] == ' ' {
		return false
	}

	// check if we have more than 2 space separation between words inside name
	if strings.Count(value, " ") != len(strings.Fields(value))-1 {
		return false
	}

	// check if name starts with a number (positive or negative)
	i := 0
	n, _ := fmt.Sscanf(value, "%d", &i)
	if n > 0 {
		// number is found at start of name
		return false
	}

	return true
}

func (k KeywordProps) ArgsValidator(value string) bool {
	// check if args starts/ends with space
	if value[0] == ' ' || value[len(value)-1] == ' ' {
		return false
	}

	return !strings.ContainsAny(value, notAllowedChars)
}

func (k KeywordProps) DocsValidator(value string) bool {
	// check if args starts/ends with space
	if value[0] == ' ' || value[len(value)-1] == ' ' {
		return false
	}

	return !strings.ContainsAny(value, notAllowedChars)
}

func (k KeywordProps) KwTypeValidator(value string) bool {
	return value == Business || value == Technical
}

type Keyword struct {
	gorm.Model
	KeywordProps
}

type User struct {
	gorm.Model
	Username string `gorm:"index;unique"`
	Email    string `gorm:"index;unique"`
	PassHash string
	Salt     string
	// role used to be enum. is gorm supports enums?
	Role string `gorm:"default:User"`
}

type History struct {
	gorm.Model
	Change string `gorm:"not null"`
}
