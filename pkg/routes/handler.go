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

type RouteHandler interface {
	Handle(*fiber.Ctx) error
}

type Route struct {
	Route    string
	Filename string
	Layouts  []string
}

type IndexRoute struct {
	Route
	db *gorm.DB
}

func NewIndexRoute(route, filename string, db *gorm.DB, layouts []string) *IndexRoute {
	return &IndexRoute{
		Route: Route{
			Route:    route,
			Filename: filename,
			Layouts:  layouts,
		},
		db: db,
	}
}

func (r *IndexRoute) Handle(c *fiber.Ctx) error {
	// Render index - start with views directory
	return c.Render(r.Filename, UpdateFiberMap(UrlMap, fiber.Map{
		"Title": "Keyword storage",
	}), r.Layouts...)
}
