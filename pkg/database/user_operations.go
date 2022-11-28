package database

import "gorm.io/gorm"

func GetUserByUsername(db *gorm.DB, username string) (User, error) {
	var user User

	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil{
		return User{}, result.Error
	}

	return user, nil
}

func GetUserById(db *gorm.DB, id int) (User, error) {
	var user User
	result := db.First(&user, id)
	if result.Error != nil{
		return User{}, result.Error
	}

	return user, nil
}