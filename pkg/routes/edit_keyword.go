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

	return c.WithUrls().Render(EditView, data)
}

func (r *Router) HandleEditKeywordPost(c *Ctx) error {
	data := c.FlashData()

	kwType := c.Params("kw_type")

	redirectUrl := RouteForKeywordType(kwType)

	kwId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.WithError(
			fmt.Sprintf("Keyword id must be number, got: %s", c.Params("id")),
		).Redirect(redirectUrl)
	}

	kwName := c.FormValue("name")
	args := c.FormValue("args")
	docs := c.FormValue("docs")

	err = database.UpdateKeyword(r.db, kwId, kwName, args, docs)

	if err != nil {
		// TODO: keep keyword data when going back
		return c.WithError(fmt.Sprintf(
			"Failed to edit Keyword '%s'!", kwName),
		).Redirect(EditKwdUrl)
	}

	data["Title"] = fmt.Sprintf("Edit %s Keyword", kwName)
	data["KwName"] = kwName
	data["Args"] = args
	data["Docs"] = docs

	return c.WithUrls().WithSuccess(fmt.Sprintf(
		"Keyword '%s' was successfully updated.", kwName),
	).Redirect(redirectUrl)
}
