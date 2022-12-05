package routes

import (
	"fmt"

	"github.com/AngelVI13/fiber-investigation/pkg/auth"
	"github.com/AngelVI13/fiber-investigation/pkg/database"
)

func (r *Router) HandleUserPanelGet(c *Ctx) error {
	username, err := auth.GetActiveUsername(c.Ctx)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"failed to get current user, error: %s", err),
		).RedirectBack(IndexUrl)
	}

	user, err := database.GetUserByUsername(r.db, username)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"Failed to find user '%s' in database, error: %s", username, err,
		)).RedirectBack(IndexUrl)
	}

	data := c.FlashData()
	data["Title"] = "User Panel"
	data["Username"] = username
	data["Email"] = user.Email
	data["Role"] = user.Role

	return c.Render(UserPanelView, data)
}

func (r *Router) HandleUserPanelPost(c *Ctx) error {
	username, err := auth.GetActiveUsername(c.Ctx)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"failed to get current user, error: %s", err),
		).RedirectBack(IndexUrl)
	}

	user, err := database.GetUserByUsername(r.db, username)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"Failed to find user '%s' in database, error: %s", username, err,
		)).RedirectBack(IndexUrl)
	}

	oldPwd := c.FormValue("old_password")
	if !database.CheckPasswordHash(user.Password, oldPwd) {
		return c.WithError("Old password is not correct").RedirectBack(IndexUrl)
	}

	newPwd := c.FormValue("new_password")
	repeatPwd := c.FormValue("repeat_password")

	if newPwd != repeatPwd {
		return c.WithError("Passwords do not match").RedirectBack(IndexUrl)
	}

	user.Password, err = database.HashPassword(newPwd)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"failed to hash password, error: %s", err),
		).RedirectBack(IndexUrl)
	}

	r.db.Save(user)
	return c.WithSuccess("Password was changed successfully").RedirectBack(IndexUrl)
}
