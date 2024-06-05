package parsers

import "newsAggr/cmd/types"

type Parser interface {
	Parse(params *types.FilteringParams) []types.News
}

func ParseWithParams(format string, params *types.FilteringParams) []types.News {
	g := GoGatorParsingFactory{}
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
