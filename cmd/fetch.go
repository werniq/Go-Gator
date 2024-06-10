package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"newsAggr/cmd/parsers"
	"newsAggr/cmd/templates"
	"newsAggr/cmd/types"
)

var fetchNews = &cobra.Command{
	Use:   "fetch",
	Short: "Fetching news from downloaded data",
	Long: "This command parses HTML, XML and JSON files sorts them by given arguments, and returns list of news" +
		"based on provided flags",

	Run: func(cmd *cobra.Command, args []string) {
		// retrieve optional parameters
		keywordsFlag, _ := cmd.Flags().GetString("keywords")
		startingTimestamp, _ := cmd.Flags().GetString("ts-from")
		endingTimestamp, _ := cmd.Flags().GetString("ts-to")
		sources, _ := cmd.Flags().GetString("sources")

		// Validate user arguments
		if err := validateDate(startingTimestamp); err != nil {
			log.Fatalln(err)
		}
		if err := validateDate(endingTimestamp); err != nil {
			log.Fatalln(err)
		}

		// Split and validate sources
		filters := types.NewFilteringParams(keywordsFlag, startingTimestamp, endingTimestamp)

		// parsing news by sources and applying params to those news
		news := parsers.ApplyParams(parsers.ParseBySource(sources), filters)

		// output using go templates
		if err := templates.ParseTemplate(filters, news); err != nil {
			panic(err)
		}

		log.Println(len(news))
	},
}

// AddFetchNewsCmd attaches fetchNews command to rootCmd
func AddFetchNewsCmd() *cobra.Command {
	fetchNews.Flags().String("keywords", "", "Topic on which news will be fetched (if empty, all news will be fetched, regardless of the theme). Separate them with ',' ")
	fetchNews.Flags().String("date-from", "", "Retrieve news based on their published date | Format 2024-05-24")
	fetchNews.Flags().String("date-end", "", "Retrieve news, where published date is not more then this value | Format 2024-05-24")
	fetchNews.Flags().String("sources", "all", "Supported sources: [abc, bbc, nbc, usatoday, washingtontimes, all]")

	return fetchNews
}
