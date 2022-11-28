package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/flash"
)

type MessageLevel string

const (
	LevelPrimary MessageLevel = "primary"
	LevelSuccess MessageLevel = "success"
	LevelWarning MessageLevel = "warning"
	LevelDanger  MessageLevel = "danger"
)

// Ctx Wraps a Ctx in order to attach utility
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

func (c *Ctx) FlashData() fiber.Map {
	return flash.Get(c.Ctx)
}

func (c *Ctx) WithInfo(message string) *Ctx {
	data := fiber.Map{
		"Message": message,
		"Level":   LevelPrimary,
	}
	flash.WithInfo(c.Ctx, data)
	return c
}

func (c *Ctx) WithSuccess(message string) *Ctx {
	data := fiber.Map{
		"Message": message,
		"Level":   LevelSuccess,
	}
	flash.WithSuccess(c.Ctx, data)
	return c
}

func (c *Ctx) WithError(message string) *Ctx {
	data := fiber.Map{
		"Message": message,
		"Level":   LevelDanger,
	}
	flash.WithError(c.Ctx, data)
	return c
}

func (c *Ctx) WithWarning(message string) *Ctx {
	data := fiber.Map{
		"Message": message,
		"Level":   LevelWarning,
	}
	flash.WithWarn(c.Ctx, data)
	return c
}
