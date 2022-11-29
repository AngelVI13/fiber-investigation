package session

import (
	"fmt"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx, user *database.User) error {
	session, err := SessionStore.Get(c)

	if err != nil {
		return fmt.Errorf("failed to get session. Error %s", err.Error())
	}

	username, err := GetActiveUsername(c)
	if err == nil {
		return fmt.Errorf("username %s is already loged in", username)
	}
	session.Set(SessionUsername, user.Username)
	session.Set(SessionRole, string(user.Role))
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

func GetActiveUsername(c *fiber.Ctx) (string, error) {
	session, err := SessionStore.Get(c)
	if err != nil {
		return "", fmt.Errorf("failed to get session. Error %s", err.Error())
	}

	username := session.Get(SessionUsername)

	usernameStr, ok := username.(string)
	if !ok {
		return "", fmt.Errorf("session variable username is not a string")
	}

	return usernameStr, nil
}

func IsAdmin(c *fiber.Ctx) bool {
	session, err := SessionStore.Get(c)
	if err != nil {
		return false
	}
	role := session.Get(SessionRole)
	return role == string(database.RoleAdmin)
}
