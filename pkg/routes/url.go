package routes

import "fmt"

const (
	IndexUrl         = "/"
	BusinessKwdsUrl  = "/business_keywords"
	TechnicalKwdsUrl = "/technical_keywords"
	AllKwdsUrl       = "/all_keywords"
	CreateKwdUrl     = "/create"
	EditKwdUrl       = "/edit"
	DeleteKwdUrl     = "/delete"
	ChangelogUrl     = "/changelog"
	ExportCsvUrl     = "/export/csv"
	ExportStubsUrl   = "/export/stubs"
	ImportCsvUrl     = "/import/csv"
	RegisterUserUrl  = "/register_user"
	LoginUrl         = "/login"
	LogoutUrl        = "/logout"
)

var (
	CreateKwdUrlFull = fmt.Sprintf("%s/:kw_type", CreateKwdUrl)
	EditKwdUrlFull   = fmt.Sprintf("%s/:id/:kw_type", EditKwdUrl)
	DeleteKwdUrlFull = fmt.Sprintf("%s/:id/:kw_type", DeleteKwdUrl)
)

var UrlMap = map[string]string{
	"IndexUrl":         IndexUrl,
	"BusinessKwdsUrl":  BusinessKwdsUrl,
	"TechnicalKwdsUrl": TechnicalKwdsUrl,
	"AllKwdsUrl":       AllKwdsUrl,
	"CreateKwdUrl":     CreateKwdUrl,
	"EditKwdUrl":       EditKwdUrl,
	"DeleteKwdUrl":     DeleteKwdUrl,
	"ChangelogUrl":     ChangelogUrl,
	"ExportCsvUrl":     ExportCsvUrl,
	"ExportStubsUrl":   ExportStubsUrl,
	"ImportCsvUrl":     ImportCsvUrl,
	"RegisterUserUrl":  RegisterUserUrl,
	"LoginUrl":         LoginUrl,
	"LogoutUrl":        LogoutUrl,
}

const (
	MainLayoutView  = "views/layouts/main"
	IndexView       = "views/index"
	KeywordsView    = "views/keywords"
	CreateView      = "views/create"
	EditView        = "views/edit"
	ChangelogView   = "views/changelog"
	ImportCsvView   = "views/import_csv"
	ExportCsvView   = "views/export_csv"
	ExportStubsView = "views/export_stubs"
	RegisterView    = "views/register"
	LoginView       = "views/login"
)
