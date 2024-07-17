package cmd

import (
	"github.com/spf13/cobra"
	"gogator/cmd/filters"
	"gogator/cmd/parsers"
	"gogator/cmd/templates"
	"gogator/cmd/types"
	"gogator/cmd/validator"
	"log"
	"strings"
)

const (
	KeywordFlag  = "keywords"
	DateFromFlag = "date-from"
	DateEndFlag  = "date-end"
	SourcesFlag  = "sources"
)

var errorMessages = map[string]string{
	"flag accessed but not defined": "Unsupported flag: ",
}

// FetchNewsCmd initializes and returns command to fetch news
// This command opens prepared files and parses their data into an array of articles.
//
// It accepts few flags: keywords, date-from, date-end, and sources.
// All of them will be used to filter retrieved news, if asked:
// Filtering by keyword will remove all articles that do not contain provided keywords. Should be separated by ','
// Date-From and Date-End are used to validate article publishing date: it will be included if it falls in range
// specified ones
// Sources flag will be defining from what sources you want to get articles from: ABC, BBC, Usa Today, Washington Times
// or all from above.
func FetchNewsCmd() *cobra.Command {
	fetchNews := &cobra.Command{}
	fetchNews.Flags().String(KeywordFlag, "", "Topic on which news will be fetched (if empty, all news will be fetched, regardless of the theme). Separate them with ',' ")
	fetchNews.Flags().String(DateFromFlag, "", "Retrieve news based on their published date | Format 2024-05-24")
	fetchNews.Flags().String(DateEndFlag, "", "Retrieve news, where published date is not more then this value | Format 2024-05-24")
	fetchNews.Flags().String(SourcesFlag, "", "Supported sources: [abc, bbc, nbc, usatoday, washingtontimes]")

	fetchNews.Use = "fetch"
	fetchNews.Short = "Fetching news from downloaded data"
	fetchNews.Long = "This command parses HTML, XML and JSON files sorts them by given arguments, and returns list of news" +
		"based on mentioned flags"

	fetchNews.Run = func(cmd *cobra.Command, args []string) {
		// retrieve optional parameters
		keywords, err := cmd.Flags().GetString(KeywordFlag)
		checkFlagErr(err)

		dateFrom, err := cmd.Flags().GetString(DateFromFlag)
		checkFlagErr(err)

		dateEnd, err := cmd.Flags().GetString(DateEndFlag)
		checkFlagErr(err)

		sources, err := cmd.Flags().GetString(SourcesFlag)
		checkFlagErr(err)

		// Validate user arguments
		if dateEnd != "" && dateFrom != "" {
			if dateFrom > dateEnd {
				log.Fatalln("Date from can not be after date end.")
			}
		}

		v := &validator.ArgValidator{}
		err = v.Validate(sources, dateFrom, dateEnd)
		if err != nil {
			log.Fatalln(err)
		}

		// Split and validate sources
		f := types.NewFilteringParams(keywords, dateFrom, dateEnd, sources)

		news, err := parsers.ParseBySource(sources)
		if err != nil {
			log.Fatalln("Error parsing news: ", err)
		}

		news = filters.Apply(news, f)

		err = templates.PrintTemplate(f, news)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return fetchNews
}

func checkFlagErr(err error) {
	if err != nil {
		for substr, msg := range errorMessages {
			if strings.Contains(err.Error(), substr) {
				log.Fatalln(msg, err)
				return
			}
		}

		log.Fatalln("Error parsing flags: ", err)
	}
}
