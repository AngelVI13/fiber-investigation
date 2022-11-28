package routes

import (
	"fmt"
	"strconv"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
)

func (r *Router) HandleDeleteKeyword(c *Ctx) error {
	kwType := c.Params("kw_type")

	redirectUrl := RouteForKeywordType(kwType)

	kwId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"Keyword id must be number, got: %s", c.Params("id")),
		).Redirect(redirectUrl)
	}

	err = database.DeleteKeyword(r.db, kwId)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"Failed to delete Keyword. Id: %d", kwId),
		).Redirect(redirectUrl)
	}

	return c.WithSuccess(fmt.Sprintf(
		"Keyword deleted successfully. Id: %d", kwId),
	).Redirect(redirectUrl)
}
