package routes

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/gocarina/gocsv"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) HandleExportCsvGet(c *fiber.Ctx) error {
	return r.renderMainLayout(c, "views/export_csv", fiber.Map{
		"Title":        "Export keywords as CSV",
		"ExportBtnTxt": "Download",
	})
}

func (r *Router) HandleExportCsvPost(c *fiber.Ctx) error {
	// TODO: abstract away database layer to something like r.db.Keywords()
	var keywords []database.Keyword

	result := r.db.Where("valid_to IS NULL").Find(&keywords)
	if result.Error != nil {
		addMessage("There are no keywords", LevelPrimary)
		// this will reload page and show message
		return r.HandleExportCsvGet(c)
	}

	filename, err := generateCsvFile("keywords.csv", keywords)
	if err != nil {
		addMessage(
			fmt.Sprintf("Error while generating csv file: %v", err),
			LevelDanger,
		)
		return r.HandleExportCsvGet(c)
	}

	c.Attachment(filepath.Base(filename))
	return c.SendFile(filename, true)
}

func generateCsvFile(filename string, keywords []database.Keyword) (string, error) {
	_ = os.Remove(filename)

	csvContents, err := generateCsv(keywords)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(filename, []byte(csvContents), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write to csv file: %v", err)
	}

	return filepath.Abs(filename)
}

func generateCsv(keywords []database.Keyword) (string, error) {
    var keywordsCsv []*database.KeywordProps

    // NOTE: need to use index in order to take a pointer of correct element
    for i := range keywords {
        keywordsCsv = append(keywordsCsv, &keywords[i].KeywordProps)
    }

    // TODO: What to do with the separator character? 
    // Can't use comma cause this might be used in the docs or args or impl
	gocsv.SetCSVWriter(func(out io.Writer) *gocsv.SafeCSVWriter {
		writer := csv.NewWriter(out)
		writer.Comma = '|'
		return gocsv.NewSafeCSVWriter(writer)
	})

	csvContents, err := gocsv.MarshalString(&keywordsCsv)
	if err != nil {
		return "", fmt.Errorf("marshalling error: %v", err)
	}
	return csvContents, nil
}

