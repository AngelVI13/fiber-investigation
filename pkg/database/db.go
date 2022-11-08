package database

import(
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
)

func Create() *gorm.DB {
    Db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("Failed to init db")
    }
    // Migrate the schema
    Db.AutoMigrate(&Keyword{}, &User{})
    return Db
}

