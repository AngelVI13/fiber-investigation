package routes

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
}

const (
	IndexView       = "views/index"
	KeywordsView    = "views/keywords"
	CreateView      = "views/create"
	EditView        = "views/edit"
	ChangelogView   = "views/changelog"
	ImportCsvView   = "views/import_csv"
	ExportCsvView   = "views/export_csv"
	ExportStubsView = "views/export_stubs"
)
