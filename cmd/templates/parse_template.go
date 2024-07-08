package templates

import (
	"fmt"
	"html/template"
	"newsaggr/cmd/types"
	"os"
	"strings"
	"time"
)

var (
	// templateFuncs are functions which can be called in the template
	templateFuncs = template.FuncMap{
		"highlight":  highlight,
		"formatDate": formatDate,
		"contains":   contains,
	}
)

const (
	BaseTemplate     = "article.plain.tmpl"
	BaseTemplatePath = "\\templates\\article.plain.tmpl"
)

func PrintTemplate(f *types.FilteringParams, articles []types.News) error {
	var err error
	cwdPath, err := os.Getwd()
	if err != nil {
		return err
	}

	sortNewsByPubDate(articles)

	tmpl := template.Must(template.New(BaseTemplate).Funcs(templateFuncs).ParseFiles(cwdPath + BaseTemplatePath))

	data := types.TemplateData{
		NewsItems:  articles,
		FilterInfo: "Applied Filters: " + fmt.Sprintf("%v", f),
		TotalItems: len(articles),
		Keywords:   strings.Split(f.Keywords, ","),
	}

	for i, v := range data.Keywords {
		if v == "" {
			data.Keywords = append(data.Keywords[:i], data.Keywords[i+1:]...)
		}
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		return err
	}

	return nil
}

// Custom function to highlight keywords
func highlight(content string, keywords []string) string {
	if keywords == nil {
		return content
	}

	for _, keyword := range keywords {
		content = strings.ReplaceAll(content, keyword, "[!]"+keyword+"[!]")
	}
	return content
}

// Custom function to format date
func formatDate(t time.Time, layout string) string {
	return t.Format(layout)
}

func contains(s string, arr []string) bool {
	for _, keyword := range arr {
		if strings.Contains(s, keyword) {
			return true
		}
	}
	return false
}
