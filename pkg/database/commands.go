package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

//InsertNewKeyword insert new keyword to database. In case it already exist, raises error
func InsertNewKeyword(db *gorm.DB, name string, args string, docs string, kw_type string) error {
	// first check that there is no kw with this name. this looks like not optimal solution
	var keyword Keyword
	kw := db.Where("name = ? AND valid_to IS NOT NULL", name).First(&keyword)
	// kw with given name exists
	if kw.Error == nil {
		return errors.New("Keyword already exist!")
	}
	keyword_record := Keyword{Name: name, Args: args, Docs: docs, KwType: kw_type}
	change := fmt.Sprintf("Add keyword '%s'", name)
	history_record := History{Change: change}
	db.Create(&keyword_record)
	db.Create(&history_record)

	return nil
}

//InsertNewKeyword update given keyword record in database.
func UpdateKeyword(db *gorm.DB, id int, name string, args string, docs string) error {
	var keyword Keyword
	kw := db.First(&keyword, id)
	// kw with given id does not exists
	if kw.Error != nil {
		return errors.New(fmt.Sprintf("Failed to get Keyword with given ID: %d", id))
	}
	now := time.Now()
	keyword.ValidTo = &now
	db.Save(&keyword)

	keyword_record := Keyword{Name: name, Args: args, Docs: docs, KwType: keyword.KwType}
	keyword_record.CreatedAt = keyword.CreatedAt
	change := fmt.Sprintf("Update keyword '%s'", name)
	history_record := History{Change: change}
	db.Create(&history_record)
	db.Create(&keyword_record)

	return nil
}

//DeleteKeyword delete keyword record by id in database.
func DeleteKeyword(db *gorm.DB, id int) error {
	var keyword Keyword
	kw := db.First(&keyword, id)
	// kw with given id does not exists
	if kw.Error != nil {
		return errors.New(fmt.Sprintf("Failed to get Keyword with given ID: %d", id))
	}

	now := time.Now()
	keyword.ValidTo = &now
	db.Save(&keyword)

	change := fmt.Sprintf("Delete keyword with ID: '%d'", id)
	history_record := History{Change: change}
	db.Create(&history_record)

	return nil
}
