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

const (
	UsaTodayKeySelector = "div.gnt_m_flm a"
	TitleTag            = "data-c-br"
	TimestampTag        = "div.gnt_m_flm_sbt"
	TimestampAttribute  = "data-c-dt"
)

// Parse function is required for HtmlParser struct, in order to implement NewsParser interface, for data formatted in html
func (hp HtmlParser) Parse() ([]types.News, error) {
	var news []types.News

	data, err := extractFileData(sourceToEndpoint[hp.Source])
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	doc.Find(UsaTodayKeySelector).Each(func(i int, selection *goquery.Selection) {
		title, _ := selection.Attr(TitleTag)
		title = strings.TrimSpace(title)
		description := strings.TrimSpace(selection.Text())
		pubDate := strings.TrimSpace(selection.Find(TimestampTag).AttrOr(TimestampAttribute, ""))

		news = append(news, types.News{
			Title:       title,
			Description: description,
			PubDate:     pubDate,
		})
	})

	return news, nil
}
