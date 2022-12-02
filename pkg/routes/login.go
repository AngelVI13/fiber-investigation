package routes

import (
	"fmt"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/AngelVI13/fiber-investigation/pkg/session"
	"github.com/gookit/validate"
)

func (r *Router) HandleLoginGet(c *Ctx) error {
	_, err := session.GetActiveUsername(c.Ctx)
	if err == nil {
		return c.WithError(
			"can't open login page: user is already logged in",
		).RedirectBack(IndexUrl)
	}

	data := c.FlashData()
	data["Title"] = "Login"
	return c.Render(LoginView, data)
}

func (r *Router) HandleLoginPost(c *Ctx) error {
	var login database.LoginForm
	err := c.BodyParser(&login)

	if err != nil {
		return c.WithError(fmt.Sprintf(
			"failed to parse request body. Error: %s", err.Error()),
		).RedirectBack(IndexUrl)
	}

	v := validate.Struct(login)
	if !v.Validate() {
		return c.WithError(fmt.Sprintf(
			"credentials validation failed. Error: %s", v.Errors),
		).RedirectBack(IndexUrl)
	}

	user, err := login.CheckLogin(r.db)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"failed to login. Error: %s", err.Error()),
		).RedirectBack(IndexUrl)
	}

	err = session.Login(c.Ctx, user)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"failed to create new seesion, error: %s", err.Error()),
		).RedirectBack(IndexUrl)
	}

	return c.WithSuccess(fmt.Sprintf(
		"User %s successfully loged in", user.Username),
	).Redirect(IndexUrl)
}

func (r *Router) HandleLogout(c *Ctx) error {
	_, err := session.GetActiveUsername(c.Ctx)
	if err != nil {
		return c.WithError(
			"can't logout: user is not logged in",
		).RedirectBack(IndexUrl)
	}

	err = session.Logout(c.Ctx)

	if err != nil {
		return c.WithError(fmt.Sprintf(
			"failed to logout, error: %s", err.Error()),
		).RedirectBack(IndexUrl)
	}

	return c.WithInfo("User logged out successfully").Redirect(LoginUrl)
}
