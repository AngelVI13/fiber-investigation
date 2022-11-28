package routes

import (
	"fmt"
	"html/template"
	"strconv"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
)

func (r Router) HandleKeywordVersion(c *Ctx) error {
	data := c.FlashData()

	versionId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.WithError(
			fmt.Sprintf("Version id must be number, got: %s", c.Params("id")),
		).Redirect(IndexUrl)
	}
	kwType := c.Params("kwType")

	kwds, err := database.KeywordsForVersion(r.db, versionId, kwType)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"Failed to fetch keywords information for version: %d", versionId),
		).Redirect(IndexUrl)
	}

	latestVersion, allVersions, err := database.LatestAndAllVersions(r.db)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"Failed to get Versioning info. Error: %v", err),
		).Redirect(IndexUrl)
	}

	data["Title"] = fmt.Sprintf("%s Keywords", kwType)
	data["Keywords"] = kwds
	data["KwType"] = template.JS(kwType)
	data["Versions"] = allVersions
	data["LatestVersion"] = latestVersion.ID
	data["SelectedVersion"] = versionId

	return c.Render(KeywordsView, data)
}

func (r *Router) HandleChangelog(c *Ctx) error {
	data := c.FlashData()
	data["Title"] = "Changelog"

	var history []database.History

	result := r.db.Find(&history)
	if result.Error != nil {
		return c.WithInfo(
			"There is no versions to display",
		).Render(ChangelogView, data)
	}

	data["History"] = history
	return c.Render(ChangelogView, data)
}
