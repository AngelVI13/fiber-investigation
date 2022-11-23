package routes

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/gocarina/gocsv"
	"github.com/sujit-baniya/flash"
)

func (r *Router) HandleImportCsvGet(c *Ctx) error {
	data := flash.Get(c.Ctx)
	data["Title"] = "Import keywords from CSV file"

	return c.WithUrls().Render("views/import_csv", data)
}

func (r *Router) HandleImportCsvPost(c *Ctx) error {
	csvContents, err := csvFileContents(c)
	if err != nil {
		return c.WithError(err.Error()).Redirect(ImportCsvUrl)
	}

	keywords, err := keywordsFromCsv(csvContents)
	if err != nil {
		return c.WithError(err.Error()).Redirect(ImportCsvUrl)
	}

	errors := insertKeywordsToDb(r, keywords)
	if errors != nil {
		msg := "During import the following errors were raised:"
		for _, e := range errors {
			msg += fmt.Sprintf("\n\t* %s", e.Error())
		}
		return c.WithWarning(msg).Redirect(ImportCsvUrl)
	}

	return c.WithSuccess(
		fmt.Sprintf("Successfully imported %d keywords", len(keywords)),
	).Redirect(ImportCsvUrl)
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

func insertKeywordsToDb(
	r *Router,
	keywordsToInsert []*database.KeywordProps,
) []error {
	var (
		allKeywords []database.Keyword
		errors      []error
	)

	_ = r.db.Where("valid_to IS NULL").Find(&allKeywords)

	var keywordMap = map[string]*database.Keyword{}
	for i := range allKeywords {
		kw := allKeywords[i]
		keywordMap[kw.Name] = &kw
	}

	for _, kw := range keywordsToInsert {
		existingKeyword, found := keywordMap[kw.Name]

		if !found {
			err := database.InsertNewKeyword(r.db, kw.Name, kw.Args, kw.Docs, kw.KwType)
			if err != nil {
				errors = append(errors, fmt.Errorf("error inserting %s: %v", kw.Name, err))
			}
			continue
		}

		err := database.UpdateKeyword(r.db, int(existingKeyword.ID), kw.Name, kw.Args, kw.Docs)
		if err != nil {
			errors = append(errors, fmt.Errorf("error updating %s: %v", kw.Name, err))
		}
	}

	return errors
}

func csvFileContents(c *Ctx) (string, error) {
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
