package routes

import (
	"fmt"
	"strconv"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
)

func (r *Router) HandleDeleteKeyword(c *Ctx) error {
	kwId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"Keyword id must be number, got: %s", c.Params("id")),
		).Redirect(IndexUrl)
	}

	err = database.DeleteKeyword(r.db, kwId)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"Failed to delete Keyword. Id: %d", kwId),
		).Redirect(IndexUrl)
	}

	return c.WithSuccess(fmt.Sprintf(
		"Keyword deleted successfully. Id: %d", kwId),
	).Redirect(IndexUrl)
}
