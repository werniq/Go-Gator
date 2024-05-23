package parsers

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"newsAggr/cmd/parsingInstructions"
	"newsAggr/cmd/types"
	"newsAggr/cmd/utils"
	"newsAggr/logger"
	"reflect"
	"strings"
)

type HtmlParser struct {
}

// Parse function is required for HtmlParser struct, in order to implement NewsParser interface, for data formatted in html
func (hp HtmlParser) Parse(params *types.ParsingParams) []types.News {
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

	for _, filename := range filenames {
		data := utils.ExtractFileData(filename)

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
	}

	factory := parsingInstructions.GoGatorInstructionFactory{}

	news = ApplyParams(news, params, factory)

	return news
}
