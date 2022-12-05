package routes

import (
	"fmt"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
)

func (r *Router) HandleIndex(c *Ctx) error {
	data := c.FlashData()
	data["Title"] = "Keyword storage"

	businesskwds, err := database.BusinessKeywords(r.db)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"error while fetching business keywords: %v", err),
		).Render(IndexView, data)
	}
	technicalkwds, err := database.TechnicalKeywords(r.db)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"error while fetching business keywords: %v", err),
		).Render(IndexView, data)
	}

	data["CountBusiness"] = len(businesskwds)
	data["CountTechnical"] = len(technicalkwds)
	data["CountAll"] = len(technicalkwds) + len(businesskwds)
	return c.Render(IndexView, data)
}
