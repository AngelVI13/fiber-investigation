package routes

import (
	"fmt"
	"log"
	"strconv"

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
}

var keywords []database.Keyword

// UpdateMap update map `n` with values from map `m`
func UpdateFiberMap[T any](m map[string]T, n fiber.Map) fiber.Map {
	for k, v := range m {
		n[k] = v
	}
	return n
}

type Router struct {
	db *gorm.DB
}

func NewRouter(db *gorm.DB) *Router {
	return &Router{
		db: db,
	}
}

func (r *Router) HandleIndex(c *fiber.Ctx) error {
	// Render index - start with views directory
	return c.Render("views/index", UpdateFiberMap(UrlMap, fiber.Map{
		"Title": "Keyword storage",
	}), "views/layouts/main")
}

func (r *Router) HandleBusinessKeywords(c *fiber.Ctx) error {
	r.db.Where("kw_type = ? AND valid_to IS NULL", "business").Find(&keywords)
	return c.Render("views/keywords", UpdateFiberMap(UrlMap, fiber.Map{
		"Title":    "Business Keywords",
		"Keywords": keywords,
		"KwType":   "business",
	}), "views/layouts/main")
}

func (r *Router) HandleTechnicalKeywords(c *fiber.Ctx) error {
	r.db.Where("kw_type = ? AND valid_to IS NULL", "technical").Find(&keywords)
	return c.Render("views/keywords", UpdateFiberMap(UrlMap, fiber.Map{
		"Title":    "Technical Keywords",
		"Keywords": keywords,
		"KwType":   "technical",
	}), "views/layouts/main")
}

func (r *Router) HandleAllKeywords(c *fiber.Ctx) error {
	r.db.Where("valid_to IS NULL").Find(&keywords)
	return c.Render("views/keywords", UpdateFiberMap(UrlMap, fiber.Map{
		"Title":    "All Keywords",
		"Keywords": keywords,
		"KwType":   "all",
	}), "views/layouts/main")
}

func (r *Router) HandleCreateKeywordGet(c *fiber.Ctx) error {
	kw_type := c.Params("kw_type")
	return c.Render("views/create", UpdateFiberMap(UrlMap, fiber.Map{
		"Title": fmt.Sprintf("Add New %s Keyword", kw_type),
	}), "views/layouts/main")
}

func (r *Router) HandleCreateKeywordPost(c *fiber.Ctx) error {
	kw_type := c.Params("kw_type")

	err := database.InsertNewKeyword(r.db, c.FormValue("name"), c.FormValue("args"), c.FormValue("docs"), kw_type)
	if err != nil {
		// this should be printed as message in html
		log.Fatalf("Failed to create new kw: %s", err)
	}
	// add message that kw was successfully added
	return c.Render("views/create", UpdateFiberMap(UrlMap, fiber.Map{
		"Title": fmt.Sprintf("Add New %s Keyword", kw_type),
	}), "views/layouts/main")
}

func (r *Router) HandleEditKeywordGet(c *fiber.Ctx) error {
	kwId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Fatalf("Failed to convert Keyword id to number: %s", err)
	}
	var keyword database.Keyword
	keywordToEdit := r.db.First(&keyword, kwId)

	if keywordToEdit.Error != nil {
		log.Fatalf("Failed to get keyword to edit(ID: %d). Error: %s", kwId, keywordToEdit.Error)
	}

	return c.Render("views/edit", UpdateFiberMap(UrlMap, fiber.Map{
		"Title":  fmt.Sprintf("Edit %s Keyword", keyword.Name),
		"KwName": keyword.Name,
		"Args":   keyword.Args,
		"Docs":   keyword.Docs,
	}), "views/layouts/main")
}

func (r *Router) HandleEditKeywordPost(c *fiber.Ctx) error {
	kwId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Fatalf("Failed to convert Keyword id to number: %s", err)
	}

	kwName := c.FormValue("name")
	args := c.FormValue("args")
	docs := c.FormValue("docs")

	err = database.UpdateKeyword(r.db, kwId, kwName, args, docs)

	if err != nil {
		// this should be printed as message in html
		log.Fatalf("Failed to update Keyword: %s", err)
	}
	// flash message that kw is updated in db

	return c.Render("views/edit", UpdateFiberMap(UrlMap, fiber.Map{
		"Title":  fmt.Sprintf("Edit %s Keyword", kwName),
		"KwName": kwName,
		"Args":   args,
		"Docs":   docs,
	}), "views/layouts/main")
}

func (r *Router) HandleDeleteKeyword(c *fiber.Ctx) error {
	kwId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Fatalf("Failed to convert Keyword id to number: %s", err)
	}

	err = database.DeleteKeyword(r.db, kwId)
	if err != nil {
		// this should be printed as message in html
		log.Fatalf("Failed to delete kw: %s", err)
	}
	// add message that kw was successfully deleted

	return r.HandleIndex(c)
}
