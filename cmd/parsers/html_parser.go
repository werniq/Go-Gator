package parsers

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"newsAggr/cmd/types"
	"newsAggr/logger"
	"reflect"
	"strings"
	"sync"
)

type HtmlParser struct {
}

// Parse function is required for HtmlParser struct, in order to implement NewsParser interface, for data formatted in html
func (hp HtmlParser) Parse(params *types.FilteringParams) []types.News {
	var news []types.News

	var filenames []string

	if !reflect.DeepEqual(params.Sources, filenames) {
		sourceToFile := map[string]string{
			"usatoday": "usa-today.html",
		}

		for _, val := range params.Sources {
			if filename, ok := sourceToFile[val]; ok {
				filenames = append(filenames, filename)
			}
		}

		if len(filenames) == 0 {
			return nil
		}
	} else {
		filenames = []string{"usa-today.html"}
	}

	var wg sync.WaitGroup

	for _, filename := range filenames {
		wg.Add(1)

		go func(filename string) {
			defer wg.Done()
			data := extractFileData(filename)

			doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
			if err != nil {
				logger.ErrorLogger.Fatalf("Unable to create new document: %v\n", err)
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
		}(filename)
	}

	wg.Wait()

	return ApplyParams(news, params)
}
