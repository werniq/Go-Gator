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
	// UsaTodayKeySelector is the CSS selector for the main content block on the USA Today page
	UsaTodayKeySelector = "a.section-helper-flex.section-helper-row.ten-column.spacer-small.p1-container"

	// TitleSelector is the CSS selector used to extract the title of an article
	TitleSelector = "div.p1-title-spacer"

	// TimestampSelector is the CSS selector used to get the article's timestamp
	TimestampSelector = "lit-timestamp"

	// TimestampAttribute is the name of the element's Attribute used to extract the publication date from the timestamp element
	TimestampAttribute = "publishdate"

	// LinkAttribute is the attribute name used to get the URL link from the element
	LinkAttribute = "href"
)

// HtmlParser is a struct implementing a Parser for HTML content from a specific source
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
		title := strings.TrimSpace(selection.Find(TitleSelector).Text())
		timestamp := strings.TrimSpace(selection.Find(TimestampSelector).AttrOr(TimestampAttribute, ""))
		link := strings.TrimSpace(selection.AttrOr(LinkAttribute, ""))
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
