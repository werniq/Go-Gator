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
func (g GoGatorParsingFactory) CreateXmlParser(source string) Parser {
	return XMLParser{
		Source: source,
	}
}

// CreateJsonParser creates an instance of JsonParser
func (g GoGatorParsingFactory) CreateJsonParser(source string) Parser {
	return JsonParser{
		Source: source,
	}
}

// CreateHtmlParser creates an instance of HtmlParser
func (g GoGatorParsingFactory) CreateHtmlParser(source string) Parser {
	return HtmlParser{
		Source: source,
	}
}
