package parsers

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"gogator/cmd/types"
	"io"
	"net/http"
	"strings"
)

const (
	UsaTodayKeySelector = "a.section-helper-flex.section-helper-row.ten-column.spacer-small.p1-container"
	TitleSelector       = "div.p1-title-spacer"
	TimestampSelector   = "lit-timestamp"
	TimestampAttribute  = "publishdate"
	LinkAttribute       = "href"
)

type HtmlParser struct {
	Source string
}

// Parse function for HtmlParser struct
func (hp HtmlParser) Parse() ([]types.News, error) {
	var news []types.News

	res, err := http.Get(sourceToEndpoint[hp.Source])
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	doc.Find(UsaTodayKeySelector).Each(func(i int, selection *goquery.Selection) {
		// Extracting the title
		title := strings.TrimSpace(selection.Find(TitleSelector).Text())

		// Extracting the timestamp
		timestamp := strings.TrimSpace(selection.Find(TimestampSelector).AttrOr(TimestampAttribute, ""))

		// Extracting the image URL (take the first URL from srcset)
		link := strings.TrimSpace(selection.AttrOr(LinkAttribute, ""))

		// Extracting description (if needed, here we just combine title and description)
		description := selection.Text()

		news = append(news, types.News{
			Title:       title,
			Description: description,
			PubDate:     timestamp,
			Publisher:   hp.Source,
			Link:        link,
		})
	})

	return news, nil
}
