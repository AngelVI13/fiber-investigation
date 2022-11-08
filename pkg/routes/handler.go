package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RouteHandler interface {
	Handle(c *fiber.Ctx, db *gorm.DB)
}

type IndexRoute struct {
	filename string
	db       *gorm.DB
}

func (r *IndexRoute) Handle(c *fiber.Ctx) {
}
