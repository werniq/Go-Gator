package templates

import (
	"fmt"
	"gogator/cmd/types"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	// templateFuncs are functions which can be called in the template
	templateFuncs = template.FuncMap{
		"highlight":  highlight,
		"formatDate": formatDate,
		"contains":   contains,
		"trim":       strings.TrimSpace,
	}

	BaseTemplatePath = filepath.Join("cmd", "templates", "templates", "article.plain.tmpl")
)

const (
	BaseTemplate = "article.plain.tmpl"
)

func PrintTemplate(f *types.FilteringParams, articles []types.News) error {
	sortNewsByPubDate(articles)

	cwdPath, err := os.Getwd()
	if err != nil {
		return err
	}

	for strings.Contains(cwdPath, "cmd") {
		cwdPath = filepath.Dir(cwdPath)
	}
	cwdPath = filepath.Join(cwdPath, BaseTemplatePath)

	tmpl := template.Must(template.New(BaseTemplate).Funcs(templateFuncs).ParseFiles(cwdPath))

	for i := 0; i <= len(articles)-1; i++ {
		articles[i] = types.News{
			Title:       strings.TrimSpace(articles[i].Title),
			Description: strings.TrimSpace(articles[i].Description),
			PubDate:     articles[i].PubDate,
			Publisher:   articles[i].PubDate,
			Link:        articles[i].Link,
		}
	}

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
