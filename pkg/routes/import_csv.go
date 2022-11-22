package routes

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/gocarina/gocsv"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) HandleImportCsvGet(c *fiber.Ctx) error {
	return r.renderMainLayout(c, "views/import_csv", fiber.Map{
		"Title": "Import keywords from CSV file",
	})
}

func (r *Router) HandleImportCsvPost(c *fiber.Ctx) error {
	csvContents, err := csvFileContents(c)
	if err != nil {
		addMessage(err.Error(), LevelDanger)
		return r.HandleImportCsvGet(c)
	}

	keywords, err := keywordsFromCsv(csvContents)
	if err != nil {
		addMessage(err.Error(), LevelDanger)
		return r.HandleImportCsvGet(c)
	}

	err = insertKeywordsToDb(keywords)
	if err != nil {
		addMessage(err.Error(), LevelDanger)
		return r.HandleImportCsvGet(c)
	}

	addMessage(
		fmt.Sprintf("Successfully imported %d keywords", len(keywords)),
		LevelSuccess,
	)
	return r.HandleImportCsvGet(c)
}

func keywordsFromCsv(csvText string) ([]*database.KeywordProps, error) {
	keywords := []*database.KeywordProps{}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '|'
		return r // Allows use pipe as delimiter
	})

	if err := gocsv.UnmarshalString(csvText, &keywords); err != nil {
		return nil, fmt.Errorf("failed to unmarshal keywords csv: %v", err)
	}
	return keywords, nil
}

func insertKeywordsToDb(keywords []*database.KeywordProps) error {
	return nil
}

func csvFileContents(c *fiber.Ctx) (string, error) {
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
