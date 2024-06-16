package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"newsaggr/cmd/parsers"
	"newsaggr/cmd/templates"
	"newsaggr/cmd/types"
	"strings"
)

func checkFlagErr(err error) {
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "flag accessed but not defined"):
			log.Fatalln("Unsupported flag: ", err)
		default:
			log.Fatalln("Error parsing flags: ", err)
		}
	}
}

// AddFetchNewsCmd attaches fetchNews command to rootCmd
func AddFetchNewsCmd() *cobra.Command {
	fetchNews := &cobra.Command{}
	fetchNews.Flags().String("keywords", "", "Topic on which news will be fetched (if empty, all news will be fetched, regardless of the theme). Separate them with ',' ")
	fetchNews.Flags().String("date-from", "", "Retrieve news based on their published date | Format 2024-05-24")
	fetchNews.Flags().String("date-end", "", "Retrieve news, where published date is not more then this value | Format 2024-05-24")
	fetchNews.Flags().String("sources", "", "Supported sources: [abc, bbc, nbc, usatoday, washingtontimes, all]")

	fetchNews.Use = "fetch"
	fetchNews.Short = "Fetching news from downloaded data"
	fetchNews.Long = "This command parses HTML, XML and JSON files sorts them by given arguments, and returns list of news" +
		"based on provided flags"

	fetchNews.Run = func(cmd *cobra.Command, args []string) {
		// retrieve optional parameters
		keywords, err := cmd.Flags().GetString("keywords")
		checkFlagErr(err)
		dateFrom, err := cmd.Flags().GetString("date-from")
		checkFlagErr(err)
		dateEnd, err := cmd.Flags().GetString("date-end")
		checkFlagErr(err)
		sources, err := cmd.Flags().GetString("sources")
		checkFlagErr(err)

		// Validate user arguments
		if dateEnd > dateFrom {
			log.Fatalln("Ending date can not be more than starting date. ")
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
		filters := types.NewFilteringParams(keywords, dateFrom, dateEnd)

		// parsing news by sources and applying params to those news
		news, err := parsers.ParseBySource(sources)
		if err != nil {
			log.Fatalln("Error parsing news: ", err)
		}

		news = parsers.ApplyParams(news, filters)

		// output using go templates
		if err = templates.PrintTemplate(filters, news); err != nil {
			log.Fatalln(err)
		}

		log.Println(len(news))
	}

	return fetchNews
}
