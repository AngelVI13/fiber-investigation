package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Create(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Migrate the schema
	db.AutoMigrate(&Keyword{}, &User{}, &History{})
	return db, nil
}
