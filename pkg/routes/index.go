package routes

func (r *Router) HandleIndex(c *Ctx) error {
	data := c.FlashData()
	data["Title"] = "Keyword storage"

	return c.WithUrls().Render(IndexView, data)
}
