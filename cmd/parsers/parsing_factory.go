package parsers

// FactoryInterface is an interface that defines methods for creating different types of parsers.
type FactoryInterface interface {
	XmlParser() Parser
	JsonParser() Parser
	HtmlParser() Parser
}

// ParsingFactory is a concrete implementation of the FactoryInterface interface.
// It provides methods to create xml, json, and html data parsers
type ParsingFactory struct{}

// XmlParser creates an instance of XMLParser
func (g ParsingFactory) XmlParser(source string) Parser {
	return XMLParser{
		Source: source,
	}
}

// JsonParser creates an instance of JsonParser
func (g ParsingFactory) JsonParser(source string) Parser {
	return JsonParser{
		Source: source,
	}
}

// HtmlParser creates an instance of HtmlParser
func (g ParsingFactory) HtmlParser(source string) Parser {
	return HtmlParser{
		Source: source,
	}
}
