package routes

import (
	"fmt"
	"strings"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
)

func (r *Router) HandleCreateKeywordGet(c *Ctx) error {
	data := c.FlashData()

	kwType := c.Params("kw_type")
	data["Title"] = fmt.Sprintf("Add New %s Keyword", kwType)

	return c.Render(CreateView, data)
}

func (r *Router) HandleCreateKeywordPost(c *Ctx) error {
	data := c.FlashData()

	kwType := c.Params("kw_type")
	data["Title"] = fmt.Sprintf("Add New %s Keyword", kwType)

	redirectUrl := RouteForKeywordType(kwType)

	nameValue := c.FormValue("name")
	argsValue := c.FormValue("args")
	docsValue := c.FormValue("docs")

	notAllowedCharset := "|"

	if strings.ContainsAny(nameValue, notAllowedCharset) ||
		strings.ContainsAny(argsValue, notAllowedCharset) ||
		strings.ContainsAny(docsValue, notAllowedCharset) {
		// TODO: how to keep the filled data after the refresh
		return c.WithError(
			fmt.Sprintf(
				`Can't create new Keyword '%s'! Some of the fields below 
                contains one or more not allowed characters(%s)`,
				nameValue,
				notAllowedCharset,
			)).RedirectBack(IndexUrl)
	}

	err := database.InsertNewKeyword(
		r.db,
		nameValue,
		argsValue,
		docsValue,
		kwType,
	)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"Failed to create new Keyword '%s'!", c.FormValue("name")),
		).RedirectBack(IndexUrl)
	}

	// add message that kw was successfully added
	return c.WithSuccess(fmt.Sprintf(
		"Added new Keyword '%s'", c.FormValue("name")),
	).Redirect(redirectUrl)
}
