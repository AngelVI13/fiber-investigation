package routes

import "github.com/gofiber/fiber/v2"

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
