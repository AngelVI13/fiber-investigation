package routes

import (
	"fmt"
	"strings"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
)

var notAllowedCharset = "|"

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
	resetQueryString(c)

	nameValue := c.FormValue("name")
	argsValue := c.FormValue("args")
	docsValue := c.FormValue("docs")

	if strings.ContainsAny(nameValue, notAllowedCharset) ||
		strings.ContainsAny(argsValue, notAllowedCharset) ||
		strings.ContainsAny(docsValue, notAllowedCharset) {
		query := makeKeywordQuery(c, nameValue, argsValue, docsValue)
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
		query := makeKeywordQuery(c, nameValue, argsValue, docsValue)
		return c.WithError(fmt.Sprintf(
			"Failed to create new Keyword '%s'!", nameValue),
		).Redirect(fmt.Sprintf("%s/%s?%s", CreateKwdUrl, kwType, query))
	}

	// add message that kw was successfully added
	return c.WithSuccess(fmt.Sprintf(
		"Added new Keyword '%s'", c.FormValue("name")),
	).Redirect(redirectUrl)
}

func makeKeywordQuery(c *Ctx, name, args, docs string) string {
	// Add query args with filled-in values so user doesn't lose
	// entered data on redirect
	c.Request().URI().QueryArgs().Add("name", name)
	c.Request().URI().QueryArgs().Add("args", args)
	c.Request().URI().QueryArgs().Add("docs", docs)
	return c.Request().URI().QueryArgs().String()
}

func resetQueryString(c *Ctx) {
	c.Request().URI().SetQueryString("")
}
