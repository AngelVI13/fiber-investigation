package database

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type KeywordProps struct {
	ValidFrom      time.Time  `gorm:"autoCreateTime;not null" csv:"-"`
	ValidTo        *time.Time `csv:"-"`
	Name           string     `gorm:"not null" csv:"Name"`
	Args           string     `gorm:"not null" csv:"Args"`
	Docs           string     `gorm:"not null" csv:"Docs"`
	KwType         string     `gorm:"not null" csv:"Type"`
	Implementation string     `csv:"Implementation"`
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
