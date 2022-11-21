package database

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

//InsertNewKeyword insert new keyword to database. In case it already exist, raises error
func InsertNewKeyword(db *gorm.DB, name string, args string, docs string, kwType string) error {
	// first check that there is no kw with this name. this looks like not optimal solution
	var keyword Keyword
	kw := db.Where("name = ? AND valid_to IS NOT NULL", name).First(&keyword)
	// kw with given name exists
	if kw.Error == nil {
		return errors.New("Keyword already exist!")
	}
	keyword_record := Keyword{Name: name, Args: args, Docs: docs, KwType: kwType}
	change := fmt.Sprintf("Add %s keyword '%s'", kwType, name)
	history_record := History{Change: change}
	db.Create(&keyword_record)
	db.Create(&history_record)

	return nil
}

//UpdateKeyword update given keyword record in database.
func UpdateKeyword(db *gorm.DB, id int, name string, args string, docs string) error {
	var keyword Keyword
	result := db.First(&keyword, id)
	if result.Error != nil {
		return errors.New(fmt.Sprintf("Failed to get Keyword with given ID: %d", id))
	}
	now := time.Now()
	keyword.ValidTo = &now
	db.Save(&keyword)

	keyword_record := Keyword{Name: name, Args: args, Docs: docs, KwType: keyword.KwType}
	keyword_record.CreatedAt = keyword.CreatedAt
	change := fmt.Sprintf("Update %s keyword '%s'", keyword.KwType, keyword.Name)
	history_record := History{Change: change}
	db.Create(&history_record)
	db.Create(&keyword_record)

	return nil
}

//DeleteKeyword delete keyword record by id in database.
func DeleteKeyword(db *gorm.DB, id int) error {
	var keyword Keyword
	result := db.First(&keyword, id)
	if result.Error != nil {
		return errors.New(fmt.Sprintf("Failed to get Keyword with given ID: %d", id))
	}

	now := time.Now()
	keyword.ValidTo = &now
	db.Save(&keyword)

	change := fmt.Sprintf("Delete %s keyword '%s'", keyword.KwType, keyword.Name)
	history_record := History{Change: change}
	db.Create(&history_record)

	return nil
}

func GetAllKeywordsForVersion(db *gorm.DB, version int, kwType string) ([]Keyword, error) {
	var keywords []Keyword
	latestVersion, err := GetLatestVersion(db)
	if err != nil {
		return nil, err
	}

	if latestVersion.ID == uint(version) {
		result := db.Where("kw_type = ? and valid_to IS NULL", kwType).Find(&keywords)
		if result.Error != nil {
			return nil, errors.New(fmt.Sprintf("Failed to get '%s' keywords for version: %d", kwType, version))
		}
		return keywords, nil
	}

	var selectedVersion History
	var nextVersion History
	result := db.First(&selectedVersion, version)
	if result.Error != nil {
		return nil, errors.New(fmt.Sprintf("Failed to get version with ID: %d.", version))
	}

	result = db.First(&nextVersion, version+1)
	if result.Error != nil {
		return nil, errors.New(fmt.Sprintf("Failed to get version with ID: %d.", version+1))
	}

	// It might be the case that no kwds are found. so no error checking is done
	_ = db.Where(
		`(
			(valid_to IS NULL AND valid_from <= ?) 
			OR 
			(valid_to IS NOT NULL AND valid_from <= ? AND valid_to >= ?)
		) 
		AND kw_type = ?`,
		selectedVersion.CreatedAt,
		selectedVersion.CreatedAt,
		nextVersion.CreatedAt,
		kwType,
	).Find(&keywords)

	return keywords, nil
}

func GetVersions(db *gorm.DB) ([]History, error) {
	var allVersions []History
	result := db.Find(&allVersions)
	if result.Error != nil {
		return nil, errors.New(fmt.Sprintf("Failed to fetch version information."))
	}
	return allVersions, nil
}

func GetLatestVersion(db *gorm.DB) (History, error) {
	var latestVersion History
	result := db.Last(&latestVersion)
	if result.Error != nil {
		return latestVersion, errors.New(fmt.Sprintf("Failed to get last version information."))
	}
	return latestVersion, nil
}
