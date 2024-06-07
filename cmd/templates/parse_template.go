package templates

import (
	"fmt"
	"html/template"
	"newsAggr/cmd/filters"
	"newsAggr/cmd/types"
	"os"
	"strings"
	"time"
)

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

func ParseTemplate(f *types.FilteringParams, articles []types.News, sortBySource bool) error {
	funcMap := template.FuncMap{
		"highlight":  highlight,
		"formatDate": formatDate,
		"contains":   contains,
	}
	sortNewsByPubDate(articles)

	d, err := os.Getwd()
	if err != nil {
		return err
	}

	tmpl := template.Must(template.New("article.plain.tmpl").Funcs(funcMap).ParseFiles(d + "\\cmd\\templates\\samples\\article.plain.tmpl"))

	data := types.TemplateData{
		NewsItems:    articles,
		FilterInfo:   "Applied Filters: " + fmt.Sprintf("%v", f),
		TotalItems:   len(articles),
		SortBySource: sortBySource,
		Keywords:     filters.SplitString(f.Keywords, ","),
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		return err
	}

	return nil
}
