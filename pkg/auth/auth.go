package auth

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
	role := CurrentUserRole(c)
	return role == database.RoleAdmin
}

func CurrentUserRole(c *fiber.Ctx) database.RoleType {
	session, err := SessionStore.Get(c)
	if err != nil {
		return database.RoleAnonymous
	}

	role := session.Get(SessionRole)
	if role == nil {
		return database.RoleAnonymous
	}

	for _, roleType := range database.AllRoles() {
		if role == string(roleType) {
			return roleType
		}
	}
	return database.RoleAnonymous
}

func RolesRequires(roles ...database.RoleType) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(roles) == 0 {
			return c.Next()
		}

		currentRole := CurrentUserRole(c)

		for _, role := range roles {
			if role == currentRole {
				return c.Next()
			}
		}
		return c.SendStatus(401)
	}
}
