package routes

import (
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
}

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
