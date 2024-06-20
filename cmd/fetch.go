package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"newsaggr/cmd/filters"
	"newsaggr/cmd/parsers"
	"newsaggr/cmd/templates"
	"newsaggr/cmd/types"
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

// AddFetchNewsCmd attaches fetchNews command to rootCmd
func AddFetchNewsCmd() *cobra.Command {
	fetchNews := &cobra.Command{}
	fetchNews.Flags().String(KeywordFlag, "", "Topic on which news will be fetched (if empty, all news will be fetched, regardless of the theme). Separate them with ',' ")
	fetchNews.Flags().String(DateFromFlag, "", "Retrieve news based on their published date | Format 2024-05-24")
	fetchNews.Flags().String(DateEndFlag, "", "Retrieve news, where published date is not more then this value | Format 2024-05-24")
	fetchNews.Flags().String(SourcesFlag, "", "Supported sources: [abc, bbc, nbc, usatoday, washingtontimes, all]")

	fetchNews.Use = "fetch"
	fetchNews.Short = "Fetching news from downloaded data"
	fetchNews.Long = "This command parses HTML, XML and JSON files sorts them by given arguments, and returns list of news" +
		"based on provided flags"

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

		err = validateDate(dateFrom)
		if err != nil {
			log.Fatalln("Error validating date: ", err)
		}

		err = validateDate(dateEnd)
		if err != nil {
			log.Fatalln("Error validating date: ", err)
		}

		err = validateSources(sources)
		if err != nil {
			log.Fatalln("Error validating sources: ", err)
		}

		// Split and validate sources
		f := types.NewFilteringParams(keywords, dateFrom, dateEnd)

		// parsing news by sources and applying params to those news
		news, err := parsers.ParseBySource(sources)
		if err != nil {
			log.Fatalln("Error parsing news: ", err)
		}

		news = filters.Apply(news, f)

		// output using go templates
		if err = templates.PrintTemplate(f, news); err != nil {
			log.Fatalln(err)
		}

		log.Println(len(news))
	}

	return fetchNews
}
