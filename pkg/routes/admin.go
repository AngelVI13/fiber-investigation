package routes

import (
	"fmt"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
)

func (r *Router) HandleAdmin(c *Ctx) error {
	data := c.FlashData()
	data["Title"] = "Admin Panel"

	var users []database.User
	result := r.db.Find(&users)
	if result.Error != nil {
		return c.WithError(fmt.Sprintf(
			"Failed to get users information from database, error: %s", result.Error,
		)).RedirectBack(IndexUrl)
	}
	data["Users"] = users
	return c.Render(AdminPanelView, data)

}

func (r *Router) HandleDeleteUser(c *Ctx) error {
	username := c.Params("username")

	err := database.DeleteUser(r.db, username)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"failed to delete User '%s', error: %s", username, err.Error()),
		).RedirectBack(IndexUrl)
	}

	return c.WithSuccess(fmt.Sprintf(
		"User '%s' deleted successfully", username),
	).RedirectBack(IndexUrl)
}

func (r *Router) HandleEditUserGet(c *Ctx) error {
	username := c.Params("username")

	user, err := database.GetUserByUsername(r.db, username)

	if err != nil {
		return c.WithError(fmt.Sprintf(
			"failed to get User '%s' from database", username),
		).RedirectBack(IndexUrl)
	}
	data := c.FlashData()
	data["Title"] = "Edit user"
	data["Username"] = user.Username
	data["Email"] = user.Email
	data["Role"] = user.Role
	data["Roles"] = [2]database.RoleType{database.RoleUser, database.RoleAdmin}

	return c.Render(EditUserView, data)
}

func (r *Router) HandleEditUserPost(c *Ctx) error {
	username := c.Params("username")
	user, err := database.GetUserByUsername(r.db, username)

	if err != nil {
		return c.WithError(fmt.Sprintf(
			"failed to get User '%s' from database", username),
		).RedirectBack(AdminPanelUrl)
	}

	var userUpdated bool = false

	usernameValue := c.FormValue("username")
	if user.Username != usernameValue {
		user.Username = usernameValue
		userUpdated = true
	}

	emailValue := c.FormValue("email")
	if user.Email != emailValue {
		user.Email = emailValue
		userUpdated = true
	}

	roleValue := c.FormValue("roles")
	if string(user.Role) != roleValue {
		user.Role = database.RoleType(roleValue)
		userUpdated = true
	}

	passwordValue := c.FormValue("new_password")
	if passwordValue != "" {
		pwd_hash, err := database.HashPassword(passwordValue)
		if err != nil {
			return c.WithError(fmt.Sprintf(
				"failed to hash new user password, error: %s", err.Error()),
			).RedirectBack(AdminPanelUrl)
		}
		user.Password = pwd_hash
		userUpdated = true
	}

	if !userUpdated {
		return c.WithWarning(fmt.Sprintf(
			"No updated for user '%s' were applied", user.Username),
		).RedirectBack(AdminPanelUrl)
	}

	result := r.db.Save(user)
	if result.Error != nil {
		return c.WithError(fmt.Sprintf(
			"failed to update user, error: %s", result.Error),
		).Redirect(AdminPanelUrl)
	}

	return c.WithSuccess(fmt.Sprintf(
		"User '%s' updated successfully", user.Username),
	).Redirect(AdminPanelUrl)
}

func (r *Router) HandleAddUserGet(c *Ctx) error {
	data := c.FlashData()
	data["Title"] = "Add New User"
	return c.Render(RegisterView, data)
}

func (r *Router) HandleAddUserPost(c *Ctx) error {
	user, err := createUser(r.db, c)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"failed to add new user, error: %s", err.Error(),
		)).RedirectBack(IndexUrl)
	}

	return c.WithSuccess(fmt.Sprintf(
		"User %s was added successfully", user.Username),
	).Redirect(IndexUrl)
}
