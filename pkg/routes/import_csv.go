package routes

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) HandleImportCsvGet(c *fiber.Ctx) error {
	return r.renderMainLayout(c, "views/import_csv", fiber.Map{
		"Title": "Import keywords from CSV file",
	})
}

func (r *Router) HandleImportCsvPost(c *fiber.Ctx) error {
	csvContents, err := getCsvFileContents(c)
	if err != nil {
		addMessage(err.Error(), LevelDanger)
		return r.HandleImportCsvGet(c)
	}

	log.Println(csvContents)

	return nil
}

func getCsvFileContents(c *fiber.Ctx) (string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return "", fmt.Errorf("error while getting multipart form: %v", err)
	}

	fileList, ok := form.File["file"]
	if !ok {
		return "", fmt.Errorf("no file object found: %v", err)
	}

	if len(fileList) < 1 {
		return "", fmt.Errorf("no file header found: %v", err)
	}

	file, err := fileList[0].Open()
	if err != nil {
		return "", fmt.Errorf("couldn't open keywords csv file: %v", err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read keywords csv file: %v", err)
	}

	return string(b), nil
}
