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

// Handler Wrapper to convert handler args to expected args by fiber and
// add url map to context.
func Handler(f func(c *Ctx) error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := &Ctx{
			Ctx: c,
		}
		return f(ctx.WithUrls().WithSession())
	}
}
