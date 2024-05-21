package cmd

type ParsingFactory interface {
	CreateXMLParser() Parsers
	CreateJSONParser() Parsers
	CreateHtmlParser() Parsers
}

type GoGatorParsingFactory struct{}

func (g GoGatorParsingFactory) CreateXmlParser() Parsers {
	return XMLParser{}
}

func (g GoGatorParsingFactory) CreateJsonParser() Parsers {
	return JSONParser{}
}

func (g GoGatorParsingFactory) CreateHtmlParser() Parsers {
	return HtmlParser{}
}
