package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gookit/validate"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RoleType string

const (
	RoleUser      RoleType = "user"
	RoleAdmin              = "admin"
	RoleAnonymous          = "anonymous"
)

func AllRoles() []RoleType {
	return []RoleType{
		RoleUser,
		RoleAdmin,
		RoleAnonymous,
	}
}

type KeywordProps struct {
	ValidFrom      time.Time  `gorm:"autoCreateTime;not null" csv:"-" validate:"-"`
	ValidTo        *time.Time `csv:"-" validate:"-"`
	Name           string     `gorm:"not null" csv:"Name" validate:"required|nameValidator"`
	Args           string     `gorm:"not null" csv:"Args" validate:"argsValidator"`
	Docs           string     `gorm:"not null" csv:"Docs" validate:"required|docsValidator"`
	KwType         string     `gorm:"not null" csv:"Type" validate:"required|kwTypeValidator"`
	Implementation string     `csv:"Implementation" validate:"-"`
}

func (k KeywordProps) String() string {
	return fmt.Sprintf("%s", k.Name)
}

type Keyword struct {
	gorm.Model
	KeywordProps
}

type History struct {
	gorm.Model
	Change string `gorm:"not null"`
}

type User struct {
	ID             uuid.UUID `json:"id" gorm:"primarykey" form:"id" validate:"-"`
	Username       string    `json:"username" gorm:"username;unique" form:"username" validate:"required"`
	Email          string    `json:"email" gorm:"email;unique" form:"email" validate:"required|email"`
	Password       string    `json:"password" gorm:"password" form:"password" validate:"required"`
	RepeatPassword string    `json:"repeat_password" form:"repeat_password" validate:"required|eq_field:password" gorm:"-"`
	Role           RoleType  `json:"role" form:"role" gorm:"default:user"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// With Messages you can custom validator error messages.
func (u User) Messages() map[string]string {
	return validate.MS{
		"required": "Field {field} is required!",
		"email":    "Invalid email format",
	}
}

func HashPassword(password string) (string, error) {
	hash_bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hash_bytes), nil
}

func CheckPasswordHash(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type LoginForm struct {
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}

func (l LoginForm) Messages() map[string]string {
	return validate.MS{
		"required": "Field {field} is required!",
	}
}

func (l LoginForm) CheckLogin(db *gorm.DB) (*User, error) {
	user, err := GetUserByUsername(db, l.Username)
	if err != nil {
		return nil, fmt.Errorf("username %s does not exist", l.Username)
	}

	if !CheckPasswordHash(user.Password, l.Password) {
		return nil, fmt.Errorf("incorrect password")
	}
	return user, nil

}
