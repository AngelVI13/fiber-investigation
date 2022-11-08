import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"gorm.io/gorm"
)

//go:embed views/*
var viewsfs embed.FS

var UrlMap = map[string]string{
	"IndexUrl":         "/",
	"BusinessKwdsUrl":  "/business_keywords",
	"TechnicalKwdsUrl": "/technical_keywords",
	"AllKwdsUrl":       "/all_keywords",
	"CreateKwdUrl":     "/create",
	"EditKwdUrl":       "/edit",
}
var keywords []database.Keyword

// how to put files in folders and then to import here?

func main1() {
	engine := html.NewFileSystem(http.FS(viewsfs), ".html")

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	db := database.Create()

	// routes
	app.Get(UrlMap["IndexUrl"], func(c *fiber.Ctx) error {
		// Render index - start with views directory
		return c.Render("views/index", UpdateFiberMap(UrlMap, fiber.Map{
			"Title": "Keyword storage",
		}), "views/layouts/main")
	})

	app.Get(UrlMap["BusinessKwdsUrl"], func(c *fiber.Ctx) error {
		// Render index - start with views directory
		db.Where("kw_type = ?", "business").Find(&keywords)
		return c.Render("views/keywords", UpdateFiberMap(UrlMap, fiber.Map{
			"Title":    "Business Keywords",
			"Keywords": keywords,
		}), "views/layouts/main")
	})

	app.Get(UrlMap["TechnicalKwdsUrl"], func(c *fiber.Ctx) error {
		db.Where("kw_type = ?", "technical").Find(&keywords)
		return c.Render("views/keywords", UpdateFiberMap(UrlMap, fiber.Map{
			"Title":    "Technical Keywords",
			"Keywords": keywords,
		}), "views/layouts/main")
	})

	app.Get(UrlMap["AllKwdsUrl"], func(c *fiber.Ctx) error {
		db.Find(&keywords)
		return c.Render("views/keywords", UpdateFiberMap(UrlMap, fiber.Map{
			"Title":    "All Keywords",
			"Keywords": keywords,
		}), "views/layouts/main")
	})

	app.Get(UrlMap["CreateKwdUrl"], func(c *fiber.Ctx) error {
		return c.Render("views/create", UpdateFiberMap(UrlMap, fiber.Map{
			"Title": "Add New Keyword",
		}), "views/layouts/main")
	})

	app.Get("/:name", indexNameHandler)

	log.Fatal(app.Listen(":3000"))
}

// UpdateMap update map `n` with values from map `m`
func UpdateFiberMap[T any](m map[string]T, n fiber.Map) fiber.Map {
	for k, v := range m {
		n[k] = v
	}
	return n
}

func indexNameHandler(c *fiber.Ctx) error {
	// Render index - start with views directory
	return c.Render("views/index", UpdateFiberMap(UrlMap, fiber.Map{
		"Title": fmt.Sprintf("Missing routes for: %s!", c.Params("name")),
	}), "views/layouts/main")
}
