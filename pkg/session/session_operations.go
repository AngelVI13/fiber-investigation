package session

import (
	"fmt"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx, user database.User) error {
	session, err := SessionStore.Get(c)

	if err != nil {
		return fmt.Errorf("failed to get session. Error %s", err.Error())
	}

	username, _ := GetActiveUser(c)
	if username != nil {
		return fmt.Errorf("username %s is already loged in", username)
	}

	session.Set(SessionUser, user.Username)
	err = session.Save()
	if err != nil {
		return fmt.Errorf("failed to add username '%s' to session. Error: %s", user.Username, err.Error())
	}

	return nil
}

func Logout(c *fiber.Ctx) error {
	session, err := SessionStore.Get(c)
	if err != nil {
		return fmt.Errorf("failed to get session. Error %s", err.Error())
	}

	err = session.Destroy()
	if err != nil {
		return fmt.Errorf("failed to destroy session. Error %s", err.Error())
	}

	return nil
}

func GetActiveUser(c *fiber.Ctx) (interface{}, error) {
	session, err := SessionStore.Get(c)
	if err != nil {
		return nil, fmt.Errorf("failed to get session. Error %s", err.Error())
	}

	return session.Get(SessionUser), nil
}
