package routes

import (
	"fmt"
	"log"
	"strings"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
)

func (r *Router) HandleCreateKeywordGet(c *Ctx) error {
	data := c.FlashData()

	kwType := c.Params("kw_type")
	data["Title"] = fmt.Sprintf("Add New %s Keyword", kwType)
	data["Name"] = c.Query("name")
	data["Args"] = c.Query("args")
	data["Docs"] = c.Query("docs")

	return c.Render(CreateView, data)
}

func (r *Router) HandleCreateKeywordPost(c *Ctx) error {
	kwType := c.Params("kw_type")

	redirectUrl := RouteForKeywordType(kwType)

	// Reset query string cause otherwise FormValue takes values from
	// query first and from multipart form second
	c.Request().URI().SetQueryString("")

	nameValue := c.FormValue("name")
	argsValue := c.FormValue("args")
	docsValue := c.FormValue("docs")
	log.Println(nameValue, argsValue, docsValue)

	notAllowedCharset := "|"

	if strings.ContainsAny(nameValue, notAllowedCharset) ||
		strings.ContainsAny(argsValue, notAllowedCharset) ||
		strings.ContainsAny(docsValue, notAllowedCharset) {
		// Add query args with filled-in values so user doesn't lose
		// entered data on redirect
		c.Request().URI().QueryArgs().Add("name", nameValue)
		c.Request().URI().QueryArgs().Add("args", argsValue)
		c.Request().URI().QueryArgs().Add("docs", docsValue)
		query := c.Request().URI().QueryArgs().String()

		return c.WithError(
			fmt.Sprintf(
				`Can't create new Keyword '%s'! Some of the fields below 
                contains one or more not allowed characters(%s)`,
				nameValue,
				notAllowedCharset,
			)).Redirect(fmt.Sprintf("%s/%s?%s", CreateKwdUrl, kwType, query))
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
