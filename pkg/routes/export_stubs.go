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
	}

	log.Println(stubType)

	// TODO: is compression needed ?
	return c.SendFile("kw_manager.go", true)
}

const TabSize = 4

var Indent = strings.Repeat(" ", TabSize)

type PyStubGenerator struct {
	filename string
	template string
}

type RfStubGenerator struct {
	filename string
	template string
}

func NewRfStubGenerator() *RfStubGenerator {
	template := `
{{.Name}}
    [Documentation]  {{.Docs}}
    [Arguments]      {{.Args}}
    Log 	         NOP

    `
	filename := "stubs.robot"

	return &RfStubGenerator{
		template: template,
		filename: filename,
	}
}

func (g *RfStubGenerator) Filename() string {
	return g.filename
}

func (g *RfStubGenerator) GenerateStubs(keywords []database.Keyword) (string, error) {
    /*
        out = "*** Keywords ***\n"
        for kw in keywords:
            name = self.generate_rf_name(kw.name)
            docs = self.generate_rf_docs(kw.docs)
            kw_txt = template.format(name=name, args=kw.args, docs=docs)
            out += kw_txt
        return out
    */
	return "", nil
}

type StubGenerator interface {
	GenerateStubs([]database.Keyword) (string, error)
	Filename() string
}

func generateStubsFile(g StubGenerator, keywords []database.Keyword) (string, error) {
	filename := g.Filename()

	_ = os.Remove(filename)

	txt, err := g.GenerateStubs(keywords)
	if err != nil {
		return "", fmt.Errorf("failed to generate stubs: %v", err)
	}

	err = os.WriteFile(filename, []byte(txt), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write to stubs file: %v", err)
	}

	return filepath.Abs(filename)
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
