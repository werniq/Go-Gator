package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"newsAggr/cmd/parsers"
	"newsAggr/cmd/types"
	"newsAggr/logger"
	"strings"
)

type FetchNewsInstruction struct {
	*cobra.Command
}

func (i FetchNewsInstruction) Execute(params *types.ParsingParams, parser parsers.Parsers) []types.News {
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

		sources := strings.Split(sourcesFlag, ",")

		// TODO: repair filtering posts with multiple parameters (e.g. keywords + --ts-from)
		// TODO: repair filtering by sources
		// TODO: add readme & instructions, add comments to functions and packages
		// TODO: add instructions as functions for filtering
		// TODO: add more tests
		// TODO: implement CommandsFactory

		// initializing parsing parameters
		parsingParams := &types.ParsingParams{
			Keywords:          keywordsFlag,
			StartingTimestamp: startingTimestamp,
			EndingTimestamp:   endingTimestamp,
			Sources:           sources,
		}

		g := parsers.GoGatorParsingFactory{}

		jsonParser := g.CreateJsonParser()
		xmlParser := g.CreateXmlParser()
		htmlParser := g.CreateHtmlParser()

		var news []types.News
		news = append(news, htmlParser.Parse(parsingParams)...)
		news = append(news, jsonParser.Parse(parsingParams)...)
		news = append(news, xmlParser.Parse(parsingParams)...)

		for _, article := range news {
			logger.InfoLogger.Println(article.Title)
			logger.InfoLogger.Println(article.Description)
			logger.InfoLogger.Println(article.PubDate)
			fmt.Println("----")
		}
		fmt.Println(len(news))
	},
}

func addFetchNewsCmd() *cobra.Command {
	fetchNews.Flags().String("keywords", "", "Topic on which news will be fetched (if empty, all news will be fetched, regardless of the theme). Separate them with ',' ")
	fetchNews.Flags().String("ts-from", "", "News starting timestamp")
	fetchNews.Flags().String("ts-to", "", "News ending timestamp")
	fetchNews.Flags().String("sources", "", "News source")

	return fetchNews
}
