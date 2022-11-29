package database

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

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
