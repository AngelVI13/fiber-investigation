package routes

import (
	"fmt"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/AngelVI13/fiber-investigation/pkg/session"
	"github.com/gookit/validate"
)

func (r *Router) HandleLoginGet(c *Ctx) error {
	data := c.FlashData()
	data["Title"] = "Login"
	return c.Render(LoginView, data)
}

func (r *Router) HandleLoginPost(c *Ctx) error {
	data := c.FlashData()
	data["Title"] = "Login"
	
	var login database.LoginForm
	err := c.BodyParser(&login)

	if err != nil {
		return c.WithError(
			fmt.Sprintf("credentials validation failed. Error: %s", err.Error()),
		).RedirectBack(IndexUrl)
	}

	v := validate.Struct(login)
	if !v.Validate() {
		return c.WithError(
			fmt.Sprintf("credentials validation failed. Error: %s", v.Errors),
		).RedirectBack(IndexUrl)
	}

	user, err := login.CheckLogin(r.db)
	if err != nil {
		return c.WithError(
			fmt.Sprintf("failed to login. Error: %s", err.Error()),
		).RedirectBack(IndexUrl)
	}

	err = session.Login(c.Ctx, user)
	if err != nil {
		return c.WithError(
			fmt.Sprintf("failed to create new seesion, error: %s", err.Error()),
		).RedirectBack(IndexUrl)
	}

	return c.WithSuccess(
		fmt.Sprintf("User %s successfully loged in", user.Username),
	).Redirect(IndexUrl)
}

func (r *Router) HandleRegisterGet(c *Ctx) error {
	data := c.FlashData()
	data["Title"] = "Register New User"
	return c.Render(RegisterView, data)
}

func (r *Router) HandleRegisterPost(c *Ctx) error {
	data := c.FlashData()
	data["Title"] = "Register New User"

	var user database.User
	err := c.BodyParser(&user)

	if err != nil {
		return c.WithError(
			fmt.Sprintf("failed to register new user. Error: %s", err.Error()),
		).RedirectBack(IndexUrl)
	}

	// validate input data
	v := validate.Struct(user)
	if !v.Validate() {
		return c.WithError(
			fmt.Sprintf("failed to register new user. Error: %s", v.Errors),
		).RedirectBack(IndexUrl)
	}

	user.HashPassword()

	result := r.db.Create(&user)

	if result.Error != nil {
		return c.WithError(
			fmt.Sprintf("failed to register new user. Error: %s", result.Error),
		).RedirectBack(IndexUrl)
	}

	defer session.Login(c.Ctx, user)

	return c.WithSuccess(
		fmt.Sprintf("User %s was added successfully", user.Username),
	).Redirect(IndexUrl)
}

func (r *Router) HandleLogout(c *Ctx) error {
	err := session.Logout(c.Ctx)

	if err != nil {
		return c.WithError(fmt.Sprintf("failed to logout, error: %s", err.Error())).RedirectBack(IndexUrl)
	}

	return c.WithInfo("User logged out successfully").Redirect(LoginUrl)
}
