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

	if p, exists := formatToParser[format]; exists {
		return p.Parse(params)
	} else {
		panic("data format " + format + " is not supported")
	}

	return nil
}

func Parse(format string, g GoGatorParsingFactory) []types.News {
	formatToParser := map[string]Parser{
		"json": g.CreateJsonParser(),
		"xml":  g.CreateXmlParser(),
		"html": g.CreateHtmlParser(),
	}

	if p, exists := formatToParser[format]; exists {
		return p.Parse(nil)
	}

	return nil
}
