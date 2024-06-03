package parsers

// ParsingFactory is an interface that defines methods for creating different types of parsers.
type ParsingFactory interface {
	CreateXMLParser() Parser
	CreateJSONParser() Parser
	CreateHtmlParser() Parser
}

// GoGatorParsingFactory is a concrete implementation of the ParsingFactory interface.
// It provides methods to create xml, json, and html data parsers
type GoGatorParsingFactory struct{}

// CreateXmlParser creates an instance of XMLParser
func (g GoGatorParsingFactory) CreateXmlParser() Parser {
	return XMLParser{}
}

// CreateJsonParser creates an instance of JsonParser
func (g GoGatorParsingFactory) CreateJsonParser() Parser {
	return JsonParser{}
}

// CreateHtmlParser creates an instance of HtmlParser
func (g GoGatorParsingFactory) CreateHtmlParser() Parser {
	return HtmlParser{}
}
