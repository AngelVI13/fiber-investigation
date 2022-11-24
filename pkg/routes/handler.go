package routes

import (
	"fmt"
	"strconv"
	"strings"

	"html/template"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/flash"
	"gorm.io/gorm"
)

const (
	IndexUrl         = "/"
	BusinessKwdsUrl  = "/business_keywords"
	TechnicalKwdsUrl = "/technical_keywords"
	AllKwdsUrl       = "/all_keywords"
	CreateKwdUrl     = "/create"
	EditKwdUrl       = "/edit"
	DeleteKwdUrl     = "/delete"
	ChangelogUrl     = "/changelog"
	ExportCsvUrl     = "/export/csv"
	ExportStubsUrl   = "/export/stubs"
	ImportCsvUrl     = "/import/csv"
)

var UrlMap = map[string]string{
	"IndexUrl":         IndexUrl,
	"BusinessKwdsUrl":  BusinessKwdsUrl,
	"TechnicalKwdsUrl": TechnicalKwdsUrl,
	"AllKwdsUrl":       AllKwdsUrl,
	"CreateKwdUrl":     CreateKwdUrl,
	"EditKwdUrl":       EditKwdUrl,
	"DeleteKwdUrl":     DeleteKwdUrl,
	"ChangelogUrl":     ChangelogUrl,
	"ExportCsvUrl":     ExportCsvUrl,
	"ExportStubsUrl":   ExportStubsUrl,
	"ImportCsvUrl":     ImportCsvUrl,
}

type Router struct {
	db *gorm.DB
}

func NewRouter(db *gorm.DB) *Router {
	return &Router{
		db: db,
	}
}

// Ctx Wraps a fiber Ctx in order to attach utility
// functions (WithUrls, WithError, etc.)
type Ctx struct {
	*fiber.Ctx
}

func (c *Ctx) WithUrls() *Ctx {
	data := fiber.Map{}

	for k, v := range UrlMap {
		data[k] = v
	}

	c.Bind(data)
	return c
}

func (r *Router) HandleIndex(c *Ctx) error {
	data := c.FlashData()
	data["Title"] = "Keyword storage"

	return c.WithUrls().Render("views/index", data)
}

func (r *Router) HandleBusinessKeywords(c *Ctx) error {
	data := c.FlashData()
	data["Title"] = "Business Keywords"

	var keywords []database.Keyword

	result := r.db.Where("kw_type = ? AND valid_to IS NULL", "business").Find(&keywords)
	if result.Error != nil {
		return c.WithUrls().WithInfo(
			"There are no business keywords to display",
		).Render("views/keywords", data)
	}

	latestVersion, allVersions, err := getLatestAndAllVersions(r.db)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"Failed to get Versioning info. Error: %v", err),
		).Redirect(IndexUrl)
	}

	data["Keywords"] = keywords
	data["KwType"] = template.JS("business")
	data["Versions"] = allVersions
	data["LatestVersion"] = latestVersion.ID
	data["SelectedVersion"] = latestVersion.ID

	return c.WithUrls().Render("views/keywords", data)
}

func (r *Router) HandleTechnicalKeywords(c *Ctx) error {
	data := c.FlashData()
	data["Title"] = "Technical Keywords"

	var keywords []database.Keyword

	result := r.db.Where("kw_type = ? AND valid_to IS NULL", "technical").Find(&keywords)
	if result.Error != nil {
		return c.WithUrls().WithInfo(
			"There are no technical keywords to display",
		).Render("views/keywords", data)
	}

	latestVersion, allVersions, err := getLatestAndAllVersions(r.db)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"Failed to get Versioning info. Error: %v", err),
		).Redirect(IndexUrl)
	}

	data["Keywords"] = keywords
	data["KwType"] = template.JS("technical")
	data["Versions"] = allVersions
	data["LatestVersion"] = latestVersion.ID
	data["SelectedVersion"] = latestVersion.ID

	return c.WithUrls().Render("views/keywords", data)
}

func (r *Router) HandleAllKeywords(c *Ctx) error {
	data := c.FlashData()
	data["Title"] = "All Keywords"

	var keywords []database.Keyword

	result := r.db.Where("valid_to IS NULL").Find(&keywords)
	if result.Error != nil {
		// TODO: What to do when for a version doesn't have keywords but i still
		// wanna go back to select older version where possibly there are keywords
		return c.WithUrls().WithInfo(
			"There are no keywords to display",
		).Render("views/keywords", data)
	}

	latestVersion, allVersions, err := getLatestAndAllVersions(r.db)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"Failed to get Versioning info. Error: %v", err),
		).Redirect(IndexUrl)
	}

	data["Keywords"] = keywords
	data["KwType"] = template.JS("all")
	data["Versions"] = allVersions
	data["LatestVersion"] = latestVersion.ID
	data["SelectedVersion"] = latestVersion.ID

	return c.WithUrls().Render("views/keywords", data)
}

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

	latestVersion, allVersions, err := getLatestAndAllVersions(r.db)
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

	return c.WithUrls().Render("views/keywords", data)
}

func (r *Router) HandleCreateKeywordGet(c *Ctx) error {
	data := c.FlashData()

	kwType := c.Params("kw_type")
	data["Title"] = fmt.Sprintf("Add New %s Keyword", kwType)

	return c.WithUrls().Render("views/create", data)
}

func (r *Router) HandleCreateKeywordPost(c *Ctx) error {
	data := c.FlashData()

	kwType := c.Params("kw_type")
	data["Title"] = fmt.Sprintf("Add New %s Keyword", kwType)

	nameValue := c.FormValue("name")
	argsValue := c.FormValue("args")
	docsValue := c.FormValue("docs")

	notAllowedCharset := "|"

	if strings.ContainsAny(nameValue, notAllowedCharset) ||
		strings.ContainsAny(argsValue, notAllowedCharset) ||
		strings.ContainsAny(docsValue, notAllowedCharset) {
		// TODO: how to keep the filled data after the refresh
		return c.WithError(
			fmt.Sprintf(
				`Can't create new Keyword '%s'!
                Some of the fields below contains one or more not allowed characters(%s)`,
				nameValue,
				notAllowedCharset,
			)).RedirectBack(IndexUrl)

	}

	err := database.InsertNewKeyword(
		r.db,
		nameValue,
		argsValue,
		docsValue,
		kwType,
	)
	if err != nil {
		return c.WithError(
			fmt.Sprintf("Failed to create new Keyword '%s'!", c.FormValue("name")),
		).RedirectBack(IndexUrl)
	}

	// add message that kw was successfully added
	return c.WithSuccess(
		fmt.Sprintf("Added new Keyword '%s'", c.FormValue("name")),
	).RedirectBack(IndexUrl)
}

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

	return c.WithUrls().Render("views/edit", data)
}

func (r *Router) HandleEditKeywordPost(c *Ctx) error {
	data := c.FlashData()

	kwId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		// TODO: Maybe this shold redirect back to where we came from and show error msg
		return c.WithError(
			fmt.Sprintf("Keyword id must be number, got: %s", c.Params("id")),
		).Redirect(IndexUrl)
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

	// TODO: add kw_type in params so that we can return to keyword page for that type
	return c.WithUrls().WithSuccess(fmt.Sprintf(
		"Keyword '%s' was successfully updated.", kwName),
	).Redirect(IndexUrl)
}

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

func (r *Router) HandleChangelog(c *Ctx) error {
	data := c.FlashData()
	data["Title"] = "Changelog"

	var history []database.History

	result := r.db.Find(&history)
	if result.Error != nil {
		return c.WithUrls().WithInfo(
			"There is no versions to display",
		).Render("views/changelog", data)
	}

	data["History"] = history
	return c.WithUrls().Render("views/changelog", data)
}

func getLatestAndAllVersions(db *gorm.DB) (database.History, []database.History, error) {
	allVersions, err := database.AllVersions(db)

	if err != nil {
		return database.History{}, nil, fmt.Errorf(
			"failed to get all Versions from db. error: %v",
			err,
		)
	}
	latestVersion, err := database.LatestVersion(db)
	if err != nil {
		return database.History{}, nil, fmt.Errorf(
			"failed to get latest Version from db. error: %v",
			err,
		)
	}
	return latestVersion, allVersions, nil
}
