package types

// TemplateData is a model which is used while working with go-templates.
//
// We will pass this model as data into go templates, to display information such as:
// /  1. NewsItems  - Array of articles, which should be displayed
// /  2. FilterInfo - Display filtering info, e.g. user inputted arguments
// /  3. TotalItems - Total amount of news
// /  4. Keywords   - Array of news. It will be used when user provided specific keywords to search articles for,
// / and using this field we will highlight these keywords.
type TemplateData struct {
	NewsItems  []News
	FilterInfo string
	TotalItems int
	Keywords   []string
}
