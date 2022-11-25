package routes

import (
	"fmt"
	"html/template"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
)

func RouteForKeywordType(kwType string) string {
	url := AllKwdsUrl
	if kwType == database.Technical {
		url = TechnicalKwdsUrl
	} else if kwType == database.Business {
		url = BusinessKwdsUrl
	}

	return url
}

func (r *Router) HandleBusinessKeywords(c *Ctx) error {
	data := c.FlashData()
	data["Title"] = "Business Keywords"

	keywords, err := database.BusinessKeywords(r.db)
	if err != nil {
		return c.WithUrls().WithError(fmt.Sprintf(
			"error while fetching business keywords: %v", err),
		).Render(KeywordsView, data)
	}

	latestVersion, allVersions, err := database.LatestAndAllVersions(r.db)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"Failed to get Versioning info. Error: %v", err),
		).Redirect(IndexUrl)
	}

	data["Keywords"] = keywords
	data["KwType"] = template.JS(database.Business)
	data["Versions"] = allVersions
	data["LatestVersion"] = latestVersion.ID
	data["SelectedVersion"] = latestVersion.ID

	return c.WithUrls().Render(KeywordsView, data)
}

func (r *Router) HandleTechnicalKeywords(c *Ctx) error {
	data := c.FlashData()
	data["Title"] = "Technical Keywords"

	keywords, err := database.TechnicalKeywords(r.db)
	if err != nil {
		return c.WithUrls().WithError(fmt.Sprintf(
			"error while fetching technical keywords: %v", err),
		).Render(KeywordsView, data)
	}

	latestVersion, allVersions, err := database.LatestAndAllVersions(r.db)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"Failed to get Versioning info. Error: %v", err),
		).Redirect(IndexUrl)
	}

	data["Keywords"] = keywords
	data["KwType"] = template.JS(database.Technical)
	data["Versions"] = allVersions
	data["LatestVersion"] = latestVersion.ID
	data["SelectedVersion"] = latestVersion.ID

	return c.WithUrls().Render(KeywordsView, data)
}

func (r *Router) HandleAllKeywords(c *Ctx) error {
	data := c.FlashData()
	data["Title"] = "All Keywords"

	keywords, err := database.AllKeywords(r.db)
	if err != nil {
		return c.WithUrls().WithError(fmt.Sprintf(
			"error while fetching all keywords: %v", err),
		).Render(KeywordsView, data)
	}

	latestVersion, allVersions, err := database.LatestAndAllVersions(r.db)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"Failed to get Versioning info. Error: %v", err),
		).Redirect(IndexUrl)
	}

	data["Keywords"] = keywords
	data["KwType"] = template.JS(database.All)
	data["Versions"] = allVersions
	data["LatestVersion"] = latestVersion.ID
	data["SelectedVersion"] = latestVersion.ID

	return c.WithUrls().Render(KeywordsView, data)
}
