package database

import "gorm.io/gorm"

func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	var user User

	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func DeleteUser(db *gorm.DB, username string) error {
	var user User
	result := db.Where("username = ?", username).Delete(&user)
	return result.Error
}
