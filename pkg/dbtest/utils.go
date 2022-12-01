package dbtest

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"gorm.io/gorm"
)

func PopulateMockKeywords(db *gorm.DB) error {
	var kwType string
	for i := 1; i <= 6; i++ {
		if i <= 3 {
			kwType = database.Business
		} else {
			kwType = database.Technical
		}
		err := database.InsertNewKeyword(
			db,
			fmt.Sprintf("Name%d", i),
			fmt.Sprintf("Args%d", i),
			fmt.Sprintf("Docs%d", i),
			kwType,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

const TestDBFileName = "db_for_test.db"

func PrepareTestDb(dbPath string) (*gorm.DB, error) {
	testDb, err := database.Create(dbPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to create test db object")
	}

	err = PopulateMockKeywords(testDb)
	if err != nil {
		return nil, fmt.Errorf("Failed to populate test db with mock data")
	}
	return testDb, nil
}

// cleanupTestDb try to remove db file created for testing
func CleanupTestDb(testDb *gorm.DB, dbPath string) func() {
	return func() {
		dbInstance, _ := testDb.DB()
		_ = dbInstance.Close()
		_ = os.Remove(dbPath)
	}
}

func NewFiberTest(t *testing.T) (*fiber.App, *gorm.DB) {
	n1 := time.Now()
	// NOTE: need to provide path to root dir so 'views' folder can be accessed
	// for testing endpoints
	path := "../../"
	viewsFs := os.DirFS(path)

	engine := html.NewFileSystem(http.FS(viewsFs), ".html")

	// TODO: this is hardcoded should be using routes.MainLayoutView
	mainLayoutView := "views/layouts/main"

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: mainLayoutView,
	})
	log.Printf("APP: %v", time.Since(n1))
	n2 := time.Now()

	dbPath := TestDBFileName
	testDb, err := PrepareTestDb(dbPath)
	if err != nil {
		t.Fatalf("couldn't create test db: %v", err)
	}

	t.Cleanup(CleanupTestDb(testDb, dbPath))

	log.Printf("DB %v", time.Since(n2))

	return app, testDb
}
