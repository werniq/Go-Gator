package templates

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed samples
var articleTemplate embed.FS

func ParseTemplate(ttr string, data interface{}) (string, error) {
	var tpl bytes.Buffer
	t, err := template.New("article.plain.tmpl").ParseFS(articleTemplate, ttr)
	if err != nil {
		return "", err
	}

	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
