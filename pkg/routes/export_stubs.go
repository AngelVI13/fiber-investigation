package routes

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
	"github.com/gofiber/fiber/v2"
)

const (
	PythonStub = "python"
	RfStub     = "rf"
)

func (r *Router) HandleExportStubsGet(c *fiber.Ctx) error {
	return r.renderMainLayout(c, "views/export_stubs", fiber.Map{
		"Title":      "Download Keywords stubs:",
		"PythonStub": PythonStub,
		"RfStub":     RfStub,
	})
}

func (r *Router) HandleExportStubsPost(c *fiber.Ctx) error {
	stubType := c.FormValue("stub_type")

	// TODO: abstract away database layer to something like r.db.Keywords()
	var keywords []database.Keyword

	result := r.db.Where("valid_to IS NULL").Find(&keywords)
	if result.Error != nil {
		addMessage("There are no keywords", LevelPrimary)
		// this will reload page and show message
		return r.HandleExportStubsGet(c)
	}

	var stubGenerator StubGenerator
	if stubType == RfStub {
		stubGenerator = &RfStubGenerator{}
	}

	filename, err := generateStubsFile(stubGenerator, keywords)
	if err != nil {
		addMessage(
			fmt.Sprintf("Error while generating stubs file: %v", err),
			LevelDanger,
		)
        return r.HandleExportStubsGet(c)
	}

    c.Attachment(filepath.Base(filename))
	return c.SendFile(filename, true)
}

const TabSize = 4

var Indent = strings.Repeat(" ", TabSize)

type PyStubGenerator struct{}

type RfStubGenerator struct{}

func (g *RfStubGenerator) Template() string {
	return `
{{.Name}}
    [Documentation]  {{.Docs}}
    [Arguments]      {{.Args}}
    Log 	         NOP

    `
}

func (g *RfStubGenerator) Filename() string {
	return "stubs.robot"
}

func (g *RfStubGenerator) Header() string {
	return "*** Keywords ***\n"
}

func (g *RfStubGenerator) TemplateProps(keyword database.Keyword) map[string]any {
	return map[string]any{
		"Name": keyword.Name,
		"Docs": keyword.Docs,
		"Args": keyword.Args,
	}
}

type StubGenerator interface {
	Filename() string
	Header() string
	TemplateProps(database.Keyword) map[string]any
	Template() string
}

func generateStubsFile(g StubGenerator, keywords []database.Keyword) (string, error) {
	filename := g.Filename()

	_ = os.Remove(filename)

	txt, err := generateStubs(g, keywords)
	if err != nil {
		return "", fmt.Errorf("failed to generate stubs: %v", err)
	}

	err = os.WriteFile(filename, []byte(txt), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write to stubs file: %v", err)
	}

	return filepath.Abs(filename)
}

func generateStubs(g StubGenerator, keywords []database.Keyword) (string, error) {
	out := g.Header()

	for _, kw := range keywords {
		props := g.TemplateProps(kw)

		kwTxt, err := formatTemplate(g.Template(), props)
		if err != nil {
			return "", err
		}
		out += kwTxt
	}

	return out, nil
}

func formatTemplate(fmt string, args map[string]interface{}) (string, error) {
	var msg bytes.Buffer

	tmpl, err := template.New("").Parse(fmt)

	if err != nil {
		return fmt, err
	}

	err = tmpl.Execute(&msg, args)
	if err != nil {
		return fmt, err
	}

	return msg.String(), nil
}
