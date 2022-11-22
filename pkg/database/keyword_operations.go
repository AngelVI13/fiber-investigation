package database

import (
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
		return fmt.Errorf("keyword already exist")
	}
	keyword_record := Keyword{
		KeywordProps: KeywordProps{
			Name:   name,
			Args:   args,
			Docs:   docs,
			KwType: kwType,
		},
	}
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
		return fmt.Errorf("failed to get keyword with given id: %d", id)
	}

	now := time.Now()
	keyword.ValidTo = &now

	db.Save(&keyword)

	keyword_record := Keyword{
		KeywordProps: KeywordProps{
			Name:   name,
			Args:   args,
			Docs:   docs,
			KwType: keyword.KwType,
		},
	}
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
		return fmt.Errorf("failed to get keyword with given id: %d", id)
	}

	now := time.Now()
	keyword.ValidTo = &now
	db.Save(&keyword)

	change := fmt.Sprintf("Delete %s keyword '%s'", keyword.KwType, keyword.Name)
	history_record := History{Change: change}
	db.Create(&history_record)

	return nil
}

func KeywordsForVersion(db *gorm.DB, version int, kwType string) ([]Keyword, error) {
	var keywords []Keyword
	latestVersion, err := LatestVersion(db)
	if err != nil {
		return nil, err
	}

	var typeFilter []string
	if kwType == "all" {
		typeFilter = append(typeFilter, "business", "technical")
	} else {
		typeFilter = append(typeFilter, kwType)
	}

	if latestVersion.ID == uint(version) {
		result := db.Where("kw_type IN ? and valid_to IS NULL", typeFilter).Find(&keywords)
		if result.Error != nil {
			return nil, fmt.Errorf("failed to get '%s' keywords for version: %d", kwType, version)
		}
		return keywords, nil
	}

	var (
		selectedVersion History
		nextVersion     History
	)
	
	result := db.First(&selectedVersion, version)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get version with id: %d", version)
	}

	result = db.First(&nextVersion, version+1)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get version with id: %d", version+1)
	}

	// It might be the case that no kwds are found. so no error checking is done
	_ = db.Where(
		`(
			(valid_to IS NULL AND valid_from <= ?)
			OR
			(valid_to IS NOT NULL AND valid_from <= ? AND valid_to >= ?)
		)
		AND kw_type IN ?`,
		selectedVersion.CreatedAt,
		selectedVersion.CreatedAt,
		nextVersion.CreatedAt,
		typeFilter,
	).Find(&keywords)

	return keywords, nil
}

func AllVersions(db *gorm.DB) ([]History, error) {
	var allVersions []History
	result := db.Find(&allVersions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch version information")
	}
	return allVersions, nil
}

func LatestVersion(db *gorm.DB) (History, error) {
	var latestVersion History
	result := db.Last(&latestVersion)
	if result.Error != nil {
		return latestVersion, fmt.Errorf("failed to get last version information")
	}
	return latestVersion, nil
}
