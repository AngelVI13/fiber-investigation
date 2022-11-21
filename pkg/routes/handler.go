package routes

import (
	"fmt"
	"strconv"
	"log"

	"html/template"
	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var UrlMap = map[string]string{
	"IndexUrl":         "/",
	"BusinessKwdsUrl":  "/business_keywords",
	"TechnicalKwdsUrl": "/technical_keywords",
	"AllKwdsUrl":       "/all_keywords",
	"CreateKwdUrl":     "/create",
	"EditKwdUrl":       "/edit",
	"DeleteKwdUrl":     "/delete",
	"ChangelogUrl":     "/changelog",
	"ExportCsvUrl":     "/export/csv",
	"ExportStubsUrl":   "/export/stubs",
}

// UpdateMap update map `n` with values from map `m`
func UpdateFiberMap[T any](m map[string]T, n fiber.Map) fiber.Map {
	for k, v := range m {
		n[k] = v
	}
	// adds alert messages. message rendering is done on every update of the page.
	n["Messages"] = getMessages()

	return n
}

type Router struct {
	db         *gorm.DB
	mainLayout string
}

func NewRouter(db *gorm.DB) *Router {
	return &Router{
		db:         db,
		mainLayout: "views/layouts/main",
	}
}

// render Wrapper for c.Render that makes sure to update
// props with UrlMap & provide main layout.
func (r *Router) renderMainLayout(
	c *fiber.Ctx,
	template string,
	props fiber.Map,
	layouts ...string,
) error {
	return c.Render(
		template,
		UpdateFiberMap(UrlMap, props),
		append(layouts, r.mainLayout)...,
	)
}

func (r *Router) HandleIndex(c *fiber.Ctx) error {
	// Render index - start with views directory
    return r.renderMainLayout(c, "views/index", fiber.Map{
		"Title": "Keyword storage",
	})
}

func (r *Router) HandleBusinessKeywords(c *fiber.Ctx) error {
	var keywords []database.Keyword
	result := r.db.Where("kw_type = ? AND valid_to IS NULL", "business").Find(&keywords)
	if result.Error != nil {
		addMessage("There is no business keywords to display", LevelPrimary)
	}
	allVersions, err := database.GetVersions(r.db)
	if err != nil {
		addMessage(fmt.Sprintf("Failed to get all Versions from db. error: %v", err), LevelDanger)
		return r.HandleIndex(c)
	}
	latestVersion, err := database.GetLatestVersion(r.db)
	if err != nil {
		addMessage(fmt.Sprintf("Failed to get latest Version from db. Error: %v", err), LevelDanger)
		return r.HandleIndex(c)
	}
    return r.renderMainLayout(c, "views/keywords", fiber.Map{
		"Title":    "Business Keywords",
		"Keywords": keywords,
		"KwType":   template.JS("business"),
		"Versions": allVersions,
		"LatestVersion": latestVersion.ID,
		"SelectedVersion": latestVersion.ID,
	})
}

func (r *Router) HandleTechnicalKeywords(c *fiber.Ctx) error {
	var keywords []database.Keyword
	result := r.db.Where("kw_type = ? AND valid_to IS NULL", "technical").Find(&keywords)
	if result.Error != nil {
		addMessage("There is no business keywords to display", LevelPrimary)
	}
	allVersions, err := database.GetVersions(r.db)
	if err != nil {
		addMessage(fmt.Sprintf("Failed to get all Versions from db. error: %v", err), LevelDanger)
		return r.HandleIndex(c)
	}
	latestVersion, err := database.GetLatestVersion(r.db)
	if err != nil {
		addMessage(fmt.Sprintf("Failed to get latest Version from db. Error: %v", err), LevelDanger)
		return r.HandleIndex(c)
	}

	return r.renderMainLayout(c, "views/keywords", fiber.Map{
		"Title":    "Technical Keywords",
		"Keywords": keywords,
		"KwType":   template.JS("technical"),
		"Versions": allVersions,
		"LatestVersion": latestVersion.ID,
		"SelectedVersion": latestVersion.ID,
	})
}

func (r *Router) HandleAllKeywords(c *fiber.Ctx) error {
	var keywords []database.Keyword
	result := r.db.Where("valid_to IS NULL").Find(&keywords)
	if result.Error != nil {
		addMessage("There is no keywords to display", LevelPrimary)
	}
	allVersions, err := database.GetVersions(r.db)
	if err != nil {
		addMessage(fmt.Sprintf("Failed to get all Versions from db. error: %v", err), LevelDanger)
		return r.HandleIndex(c)
	}
	latestVersion, err := database.GetLatestVersion(r.db)
	if err != nil {
		addMessage(fmt.Sprintf("Failed to get latest Version from db. Error: %v", err), LevelDanger)
		return r.HandleIndex(c)
	}
	return r.renderMainLayout(c, "views/keywords", fiber.Map{
		"Title":    "All Keywords",
		"Keywords": keywords,
		"KwType":   template.JS("all"),
		"Versions": allVersions,
		"LatestVersion": latestVersion.ID,
		"SelectedVersion": latestVersion.ID,
	})
}

func (r *Router) HandleCreateKeywordGet(c *fiber.Ctx) error {
	kw_type := c.Params("kw_type")
	return r.renderMainLayout(c, "views/create", fiber.Map{
		"Title": fmt.Sprintf("Add New %s Keyword", kw_type),
	})
}

func (r *Router) HandleCreateKeywordPost(c *fiber.Ctx) error {
	kw_type := c.Params("kw_type")

	err := database.InsertNewKeyword(
		r.db,
		c.FormValue("name"),
		c.FormValue("args"),
		c.FormValue("docs"),
		kw_type,
	)
	if err != nil {
		addMessage(fmt.Sprintf("Failed to create new Keyword '%s'!", c.FormValue("name")), LevelDanger)
	} else {
		addMessage(fmt.Sprintf("Added new Keyword '%s'", c.FormValue("name")), LevelSuccess)
	}
	// add message that kw was successfully added
	return r.renderMainLayout(c, "views/create", fiber.Map{
		"Title": fmt.Sprintf("Add New %s Keyword", kw_type),
	})
}

func (r *Router) HandleEditKeywordGet(c *fiber.Ctx) error {
	kwId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		addMessage(fmt.Sprintf("Keyword id must be number, got: %s", c.Params("id")), LevelDanger)
		return r.HandleIndex(c)
	}
	var keyword database.Keyword
	result := r.db.First(&keyword, kwId)

	if result.Error != nil {
		addMessage(fmt.Sprintf("Failed to get Keyword (ID: %d) to edit!", kwId), LevelDanger)
		return r.HandleIndex(c)
	}

	return r.renderMainLayout(c, "views/edit", fiber.Map{
		"Title":  fmt.Sprintf("Edit %s Keyword", keyword.Name),
		"KwName": keyword.Name,
		"Args":   keyword.Args,
		"Docs":   keyword.Docs,
	})
}

func (r *Router) HandleEditKeywordPost(c *fiber.Ctx) error {
	kwId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		addMessage(fmt.Sprintf("Keyword id must be number, got: %s", c.Params("id")), LevelDanger)
		return r.HandleIndex(c)
	}

	kwName := c.FormValue("name")
	args := c.FormValue("args")
	docs := c.FormValue("docs")

	err = database.UpdateKeyword(r.db, kwId, kwName, args, docs)

	if err != nil {
		addMessage(fmt.Sprintf("Failed to edit Keyword '%s'!", kwName), LevelDanger)
	} else {
		addMessage(fmt.Sprintf("Keyword '%s' was successfully updated.", kwName), LevelSuccess)
	}

	return r.renderMainLayout(c, "views/edit", fiber.Map{
		"Title":  fmt.Sprintf("Edit %s Keyword", kwName),
		"KwName": kwName,
		"Args":   args,
		"Docs":   docs,
	})
}

func (r *Router) HandleDeleteKeyword(c *fiber.Ctx) error {
	kwId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		addMessage(fmt.Sprintf("Keyword id must be number, got: %s", c.Params("id")), LevelDanger)
		return r.HandleIndex(c)
	}

	err = database.DeleteKeyword(r.db, kwId)
	if err != nil {
		addMessage(fmt.Sprintf("Failed to delete Keyword. Id: %d", kwId), LevelDanger)
	} else {
		addMessage(fmt.Sprintf("Keyword deleted successfully. Id: %d", kwId), LevelPrimary)
	}

	return r.HandleIndex(c)
}

func (r *Router) HandleChangelog(c *fiber.Ctx) error {
	var history []database.History
	result := r.db.Find(&history)
	if result.Error != nil {
		addMessage("There is no versions to display", LevelPrimary)
	}

	return r.renderMainLayout(c, "views/changelog", fiber.Map{
		"Title":   "Changelog",
		"History": history,
	})
}

func (r Router) HandleKeywordVersion(c *fiber.Ctx) error {
	versionId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		addMessage(fmt.Sprintf("Version id must be number, got: %s", c.Params("id")), LevelDanger)
		return r.HandleIndex(c)
	}
	kwType := c.Params("kwType")

	kwds, err := database.GetAllKeywordsForVersion(r.db, versionId, kwType)
	if err != nil {
		addMessage(fmt.Sprintf("Failed to fetch keywors information for version: %d", versionId), LevelDanger)
		return r.HandleIndex(c)
	}
	allVersions, err := database.GetVersions(r.db)
	if err != nil {
		addMessage(fmt.Sprintf("Failed to get all Versions from db. error: %v", err), LevelDanger)
		return r.HandleIndex(c)
	}
	latestVersion, err := database.GetLatestVersion(r.db)
	if err != nil {
		addMessage(fmt.Sprintf("Failed to get latest Version from db. Error: %v", err), LevelDanger)
		return r.HandleIndex(c)
	}


	return r.renderMainLayout(c, "views/keywords", fiber.Map{
		"Title":    fmt.Sprintf("%s Keywords", kwType),
		"Keywords": kwds,
		"KwType":   template.JS(kwType),
		"Versions": allVersions,
		"LatestVersion": latestVersion.ID,
		"SelectedVersion": versionId,
	})
}

func (r *Router) HandleExportCsv(c *fiber.Ctx) error {
	log.Println("Export CSV")
	return nil
}

