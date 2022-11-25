package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

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

	return c.WithUrls().Render(IndexView, data)
}
