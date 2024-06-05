package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"log"
	"newsAggr/cmd/filters"
	"newsAggr/cmd/parsers"
	"newsAggr/cmd/types"
	"newsAggr/logger"
	"testing"
)

func Test_addFetchNewsCmd(t *testing.T) {
	testFetchNews := &cobra.Command{
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

			news := parsers.ParseWithParams("json", filteringParams)
			news = append(news, parsers.ParseWithParams("xml", filteringParams)...)
			news = append(news, parsers.ParseWithParams("html", filteringParams)...)

			for _, article := range news {
				logger.InfoLogger.Println(article.Title)
				logger.InfoLogger.Println(article.Description)
				logger.InfoLogger.Println(article.PubDate)
				logger.InfoLogger.Println(article.Link)
				fmt.Println("----")
			}
			fmt.Println("Articles retrieved: ", len(news))
		},
	}

	assert.Equal(t, testFetchNews.Short, fetchNews.Short)
}
