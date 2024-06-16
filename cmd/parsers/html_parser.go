package parsers

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"newsaggr/cmd/types"
	"strings"
)

type HtmlParser struct {
	Source string
}

// Parse function is required for HtmlParser struct, in order to implement NewsParser interface, for data formatted in html
func (hp HtmlParser) Parse() ([]types.News, error) {
	var news []types.News

	data, err := extractFileData(sourceToFile[hp.Source])
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	doc.Find("div.gnt_m_flm a").Each(func(i int, selection *goquery.Selection) {
		title, _ := selection.Attr("data-c-br")
		title = strings.TrimSpace(title)
		description := strings.TrimSpace(selection.Text())
		pubDate := strings.TrimSpace(selection.Find("div.gnt_m_flm_sbt").AttrOr("data-c-dt", ""))

		news = append(news, types.News{
			Title:       title,
			Description: description,
			PubDate:     pubDate,
		})
	})

	return news, nil
}
