package database

import (
	"time"

	"gorm.io/gorm"
)

type Keyword struct {
	gorm.Model
	ValidFrom      time.Time `gorm:"autoCreateTime;not null"`
	ValidTo        *time.Time
	Name           string `gorm:"not null,unique"`
	Args           string `gorm:"not null"`
	Docs           string `gorm:"not null"`
	KwType         string `gorm:"not null"`
	Implementation string
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
