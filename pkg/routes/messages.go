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

func (c *Ctx) FlashData() fiber.Map {
	return flash.Get(c.Ctx)
}

func (c *Ctx) WithInfo(message string) *fiber.Ctx {
	data := fiber.Map{
		"Message": message,
		"Level":   LevelPrimary,
	}
	return flash.WithInfo(c.Ctx, data)
}

func (c *Ctx) WithSuccess(message string) *fiber.Ctx {
	data := fiber.Map{
		"Message": message,
		"Level":   LevelSuccess,
	}
	return flash.WithSuccess(c.Ctx, data)
}

func (c *Ctx) WithError(message string) *fiber.Ctx {
	data := fiber.Map{
		"Message": message,
		"Level":   LevelDanger,
	}
	return flash.WithError(c.Ctx, data)
}

func (c *Ctx) WithWarning(message string) *fiber.Ctx {
	data := fiber.Map{
		"Message": message,
		"Level":   LevelWarning,
	}
	return flash.WithWarn(c.Ctx, data)
}
