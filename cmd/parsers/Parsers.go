package parsers

import "newsAggr/cmd/types"

type Parser interface {
	Parse(params *types.FilteringParams) []types.News
}

func ParseWithParams(format string, params *types.FilteringParams, g GoGatorParsingFactory) []types.News {
	formatToParser := map[string]Parser{
		"json": g.CreateJsonParser(),
		"xml":  g.CreateXmlParser(),
		"html": g.CreateHtmlParser(),
	}

	return formatToParser[format].Parse(params)
}

func Parse(format string, g GoGatorParsingFactory) []types.News {
	formatToParser := map[string]Parser{
		"json": g.CreateJsonParser(),
		"xml":  g.CreateXmlParser(),
		"html": g.CreateHtmlParser(),
	}

	return formatToParser[format].Parse(nil)
}
