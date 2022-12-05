package routes

import (
	"fmt"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/AngelVI13/fiber-investigation/pkg/auth"
	"github.com/google/uuid"
	"github.com/gookit/validate"
	"gorm.io/gorm"
)

// createUser parse context body and creates new User record in db  
func createUser(db *gorm.DB, c *Ctx) (*database.User, error) {
	var user database.User
	err := c.BodyParser(&user)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse request body. Error: %s", err.Error(),
		)
	}

	// validate input data
	v := validate.Struct(user)
	if !v.Validate() {
		return nil, fmt.Errorf(
			"failed to register new user. Error: %s", v.Errors,
		)
	}

	pwdHash, err := database.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to hash password. Error: %s", err.Error(),
		)
	}

	user.Password = pwdHash
	user.ID = uuid.New()

	result := db.Create(&user)

	if result.Error != nil {
		return nil, fmt.Errorf(
			"failed to register new user. Error: %s", result.Error,
		)
	}

	return &user, nil
}

func (r *Router) HandleRegisterGet(c *Ctx) error {
	data := c.FlashData()
	data["Title"] = "Register New User"
	return c.Render(RegisterView, data)
}

func (r *Router) HandleRegisterPost(c *Ctx) error {
	user, err := createUser(r.db, c)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"failed to create new user, error: %s", err.Error(),
		)).RedirectBack(IndexUrl)
	}

	err = auth.Login(c.Ctx, user)
	if err != nil {
		return c.WithWarning(fmt.Sprintf(
			"user registered successfully, but failed to login. error: %s", err),
		).RedirectBack(IndexUrl)
	}

	return c.WithSuccess(fmt.Sprintf(
		"User %s was added successfully", user.Username),
	).Redirect(IndexUrl)
}
