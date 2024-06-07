package types

type TemplateData struct {
	NewsItems    []News
	FilterInfo   string
	TotalItems   int
	SortBySource bool
	Keywords     []string
}
