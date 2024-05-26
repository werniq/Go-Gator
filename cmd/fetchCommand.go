package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"newsAggr/cmd/FilteringInstructions"
	"newsAggr/cmd/parsers"
	"newsAggr/cmd/types"
	"newsAggr/logger"
)

type FetchNewsInstruction struct {
	*cobra.Command
}

func (i FetchNewsInstruction) Execute(params *types.FilteringParams, parser parsers.Parsers) []types.News {
	return parser.Parse(params)
}

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

		sources := FilteringInstructions.Split(sourcesFlag, ",")

		// initializing parsing parameters
		parsingParams := &types.FilteringParams{
			Keywords:          keywordsFlag,
			StartingTimestamp: startingTimestamp,
			EndingTimestamp:   endingTimestamp,
			Sources:           sources,
		}

		g := parsers.GoGatorParsingFactory{}

		parsers := []parsers.Parsers{
			g.CreateJsonParser(),
			g.CreateXmlParser(),
			g.CreateHtmlParser(),
		}

		var news []types.News

		for _, parser := range parsers {
			parsedNews := parser.Parse(parsingParams)
			news = append(news, parsedNews...)
		}

		for _, article := range news {
			logger.InfoLogger.Println(article.Title)
			logger.InfoLogger.Println(article.Description)
			logger.InfoLogger.Println(article.PubDate)
			fmt.Println("----")
		}
		fmt.Println("Articles retrieved: ", len(news))
	},
}

func addFetchNewsCmd() *cobra.Command {
	fetchNews.Flags().String("keywords", "", "Topic on which news will be fetched (if empty, all news will be fetched, regardless of the theme). Separate them with ',' ")
	fetchNews.Flags().String("ts-from", "", "News starting timestamp")
	fetchNews.Flags().String("ts-to", "", "News ending timestamp")
	fetchNews.Flags().String("sources", "", "Supported sources: [abcnews, bbc, nbc, usatoday, washingtontimes]")

	return fetchNews
}
