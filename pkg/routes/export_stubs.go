package routes

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/AngelVI13/fiber-investigation/pkg/database"
)

const (
	PythonStub = "python"
	RfStub     = "rf"
)

func (r *Router) HandleExportStubsGet(c *Ctx) error {
	data := c.FlashData()

	data["Title"] = "Download Keywords stubs:"
	data["PythonStub"] = PythonStub
	data["RfStub"] = RfStub

	return c.WithUrls().Render(ExportStubsView, data)
}

func (r *Router) HandleExportStubsPost(c *Ctx) error {
	stubType := c.FormValue("stub_type")

	keywords, err := database.AllKeywords(r.db)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"error while fetching all keywords: %v", err),
		).Redirect(ExportStubsUrl)
	}

	var stubGenerator StubGenerator
	if stubType == RfStub {
		stubGenerator = &RfStubGenerator{}
	} else if stubType == PythonStub {
		stubGenerator = &PyStubGenerator{}
	} else {
		return c.WithError(fmt.Sprintf(
			"Unsupported stub type %s", strconv.Quote(stubType),
		)).Redirect(ExportStubsUrl)
	}

	filename, err := generateStubsFile(stubGenerator, keywords)
	if err != nil {
		return c.WithError(fmt.Sprintf(
			"Error while generating stubs file: %v", err),
		).Redirect(ExportStubsUrl)
	}

	c.Attachment(filepath.Base(filename))
	return c.SendFile(filename, true)
}

const TabSize = 4

var Indent = strings.Repeat(" ", TabSize)

type PyStubGenerator struct{}

func (g *PyStubGenerator) Template() string {
	return `
@keyword("{{.RawName}}")
def {{.Name}}({{.Args}}):
    """
    {{.Docs}}
    """
    pass

    `
}

func (g *PyStubGenerator) Filename() string {
	return "stubs.py"
}

func (g *PyStubGenerator) Header() string {
	return "from robot.api.deco import keyword\n"
}

func (g *PyStubGenerator) RawName(keyword string) string {
	fields := strings.Fields(keyword)
	return strings.Join(fields, " ")
}

// Name Clean keyword name to be suitable for python method name
// Remove all leading/trailing whitespaces and any extra spaces between words.
// Lowercase all characters and join words with underscores.
func (g *PyStubGenerator) Name(keyword string) string {
	fields := strings.Fields(keyword)

	var newFields []string
	for _, f := range fields {
		newFields = append(newFields, strings.ToLower(f))
	}

	return strings.Join(newFields, "_")
}

// Docs Clean keyword docs to be suitable for python docstring style
func (g *PyStubGenerator) Docs(docs string) string {
	docs = strings.TrimSpace(docs)
	lines := strings.SplitAfter(docs, "\n")

	cleanDocs := ""
	for i, line := range lines {
		if i == 0 {
			// The first line does not need the above prefix because it gets
			// added to the [Documentation] part of the docstring
			cleanDocs += line
			continue
		}

		// Add an appropriate indentation to each line of the docs
		cleanDocs += fmt.Sprintf("%s%s", Indent, line)
	}

	return cleanDocs
}

func (g *PyStubGenerator) TemplateProps(keyword database.Keyword) map[string]any {
	return map[string]any{
		"RawName": g.RawName(keyword.Name),
		"Name":    g.Name(keyword.Name),
		"Docs":    g.Docs(keyword.Docs),
		"Args":    keyword.Args,
	}
}

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

// Name Clean keyword name to be suitable for robot file style
// Remove all leading/trailing whitespaces and any extra spaces between words
func (g *RfStubGenerator) Name(keyword string) string {
	fields := strings.Fields(keyword)
	return strings.Join(fields, " ")
}

// Docs Clean keyword docs to be suitable for robot file style
func (g *RfStubGenerator) Docs(docs string) string {
	docs = strings.TrimSpace(docs)
	lines := strings.SplitAfter(docs, "\n")

	cleanDocs := ""
	for i, line := range lines {
		if i == 0 {
			// The first line does not need the above prefix because it gets
			// added to the [Documentation] part of the docstring
			cleanDocs += line
			continue
		}

		// Add '...' and an appropriate indentation to each line of the docs
		cleanDocs += fmt.Sprintf("%s...%s%s", Indent, Indent, line)
	}

	return cleanDocs
}

func (g *RfStubGenerator) TemplateProps(keyword database.Keyword) map[string]any {
	return map[string]any{
		"Name": g.Name(keyword.Name),
		"Docs": g.Docs(keyword.Docs),
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
