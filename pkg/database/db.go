package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Create(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	// Migrate the schema
	db.AutoMigrate(&Keyword{}, &User{}, &History{})
	return db, nil
}
