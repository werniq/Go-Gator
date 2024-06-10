package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"newsAggr/cmd/filters"
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
		sourcesFlag, _ := cmd.Flags().GetString("sources")

		if err := validateDate(startingTimestamp); err != nil {
			log.Fatalln(err)
		}
		if err := validateDate(endingTimestamp); err != nil {
			log.Fatalln(err)
		}

		// Split and validate sources
		sources := filters.SplitString(sourcesFlag, ",")
		if err := validateSources(sources); err != nil {
			log.Fatalln(err)
		}

		filteringParams := types.NewParams(keywordsFlag, startingTimestamp, endingTimestamp, sources)

		news := parsers.ParseBySources("json", filteringParams)
		news = append(news, parsers.ParseBySources("xml", filteringParams)...)
		news = append(news, parsers.ParseBySources("html", filteringParams)...)
		news = parsers.ApplyParams(news, filteringParams)

		if err := templates.ParseTemplate(filteringParams, news); err != nil {
			panic(err)
		}
	},
}

func AddFetchNewsCmd() *cobra.Command {
	fetchNews.Flags().String("keywords", "", "Topic on which news will be fetched (if empty, all news will be fetched, regardless of the theme). Separate them with ',' ")
	fetchNews.Flags().String("ts-from", "", "Retrieve news based on their published date | Format 2024-05-24")
	fetchNews.Flags().String("ts-to", "", "Retrieve news, where published date is not more then this value | Format 2024-05-24")
	fetchNews.Flags().String("sources", "", "Supported sources: [abc, bbc, nbc, usatoday, washingtontimes]")

	return fetchNews
}
