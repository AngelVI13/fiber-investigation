package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"

	"gorm.io/gorm"
    "gorm.io/driver/sqlite"
)

//go:embed views/*
var viewsfs embed.FS

var UrlMap = map[string]string{
	"IndexUrl": "/",
	"BusinessKwdsUrl": "/business_keywords",
	"TechnicalKwdsUrl": "/technical_keywords",
	"AllKwdsUrl": "/all_keywords",
}
var keywords []Keyword


// how to put files in folders and then to import here?
type Keyword struct {
    gorm.Model
    ValidFrom       time.Time       `gorm:"autoCreateTime;not null"`
    ValidTo         *time.Time
    Name            string          `gorm:"not null,unique"`
    Args            string          `gorm:"not null"`
    Docs            string          `gorm:"not null"`
    KwType          string          `gorm:"not null"`
    Implementation  string
}

type User struct{
    Username        string          `gorm:"index;unique"`
    Email           string          `gorm:"index;unique"`
    PassHash        string
    Salt            string
    // role used to be enum. is gorm supports enums?
    Role            string          `gorm:"default:User"`
}

func main() {
	engine := html.NewFileSystem(http.FS(viewsfs), ".html")

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

    // init db
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // Migrate the schema
    db.AutoMigrate(&Keyword{}, &User{})

    // routes
	app.Get(UrlMap["IndexUrl"], func(c *fiber.Ctx) error {
		// Render index - start with views directory
		return c.Render("views/index", UpdateFiberMap(UrlMap, fiber.Map{
			"Title":     "Keyword storage",
		}), "views/layouts/main")
	})

	app.Get(UrlMap["BusinessKwdsUrl"], func(c *fiber.Ctx) error {
		// Render index - start with views directory
		db.Where("kw_type = ?", "business").Find(&keywords)
		return c.Render("views/keywords", UpdateFiberMap(UrlMap, fiber.Map{
			"Title":     "Business Keywords",
			"Keywords":  keywords,
		}), "views/layouts/main")
	})

    app.Get(UrlMap["TechnicalKwdsUrl"], func(c *fiber.Ctx) error {
		// Render index - start with views directory
		db.Where("kw_type = ?", "technical").Find(&keywords)
		return c.Render("views/keywords", UpdateFiberMap(UrlMap, fiber.Map{
			"Title":     "Technical Keywords",
			"Keywords":  keywords,
		}), "views/layouts/main")
	})

    app.Get(UrlMap["AllKwdsUrl"], func(c *fiber.Ctx) error {
		// Render index - start with views directory
		db.Find(&keywords)
		return c.Render("views/keywords", UpdateFiberMap(UrlMap, fiber.Map{
			"Title":     "All Keywords",
			"Keywords":  keywords,
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
		"Title":     fmt.Sprintf("Missing routes for: %s!", c.Params("name")),
	}), "views/layouts/main")
}
