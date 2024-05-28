package parsers

import "newsAggr/cmd/types"

type Parsers interface {
	Parse(params *types.FilteringParams) []types.News
}

func Parse(format string, params *types.FilteringParams, g GoGatorParsingFactory) []types.News {
	var p Parsers
	switch format {
	case "json":
		p = g.CreateJsonParser()
	case "xml":
		p = g.CreateXmlParser()
	case "html":
		p = g.CreateHtmlParser()
	}

	return p.Parse(params)
}
