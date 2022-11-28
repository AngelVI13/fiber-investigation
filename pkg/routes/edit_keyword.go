package routes

import (
	"fmt"
	"strconv"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
)

func (r *Router) HandleEditKeywordGet(c *Ctx) error {
	data := c.FlashData()

	kwId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"Keyword id must be number, got: %s", c.Params("id")),
		).Redirect(IndexUrl)
	}
	var keyword database.Keyword
	result := r.db.First(&keyword, kwId)

	if result.Error != nil {
		return c.WithError(fmt.Sprintf(
			"Failed to get Keyword (ID: %d) to edit!", kwId),
		).Redirect(IndexUrl)
	}

	data["Title"] = fmt.Sprintf("Edit %s Keyword", keyword.Name)
	data["KwName"] = keyword.Name
	data["Args"] = keyword.Args
	data["Docs"] = keyword.Docs

	// If query contains name/args/docs update shown values with those from
	// query instead of their keyword version. This is done so on error
	// the user input is kept and not erased.
	queryName := c.Query("name")
	queryArgs := c.Query("args")
	queryDocs := c.Query("docs")

	if queryName != "" {
		data["KwName"] = queryName
	}
	if queryArgs != "" {
		data["Args"] = queryArgs
	}
	if queryDocs != "" {
		data["Docs"] = queryDocs
	}
	return c.Render(EditView, data)
}

func (r *Router) HandleEditKeywordPost(c *Ctx) error {
	resetQueryString(c)

	kwType := c.Params("kw_type")

	redirectUrl := RouteForKeywordType(kwType)

	kwId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.WithError(
			fmt.Sprintf("Keyword id must be number, got: %s", c.Params("id")),
		).Redirect(redirectUrl)
	}

	nameValue := c.FormValue("name")
	argsValue := c.FormValue("args")
	docsValue := c.FormValue("docs")

	err = database.UpdateKeyword(r.db, kwId, nameValue, argsValue, docsValue)

	if err != nil {
		// keep entered data inside query string
		query := makeKeywordQuery(c, nameValue, argsValue, docsValue)

		editUrl, err := editUrl(kwId, kwType, query)
		if err != nil {
			return c.WithError(fmt.Sprintf(
				`error while trying to redirect back to %s 
                page after error: couldn't format url`, EditKwdUrl),
			).Redirect(IndexUrl)
		}

		return c.WithError(fmt.Sprintf(
			"Failed to edit Keyword '%s'!", nameValue),
		).Redirect(editUrl)
	}

	return c.WithSuccess(fmt.Sprintf(
		"Keyword '%s' was successfully updated.", nameValue),
	).Redirect(redirectUrl)
}

func editUrl(kwId int, kwType, query string) (string, error) {
	url := "{{.Base}}/{{.Id}}/{{.KwType}}?{{.Query}}"

	props := map[string]any{
		"Base":   EditKwdUrl,
		"Id":     kwId,
		"KwType": kwType,
		"Query":  query,
	}
	return formatTemplate(url, props)
}
