package parsers

func determineParser(format, source string) Parser {
	switch format {
	case "json":
		return JsonParser{Source: source}
	case "xml":
		return XMLParser{Source: source}
	case "html":
		return HtmlParser{Source: source}
	}
	return nil
}

func AddNewSource(format, source, endpoint string) {
	sourceToEndpoint[source] = endpoint
	sourceToParser[source] = determineParser(format, source)
}

func GetAllSources() map[string]string {
	return sourceToEndpoint
}

func UpdateSourceEndpoint(source, newEndpoint string) {
	sourceToEndpoint[source] = newEndpoint
}

func UpdateSourceFormat(source, format string) {
	sourceToParser[source] = determineParser(format, source)
}

func DeleteSource(source string) {
	if _, exists := sourceToEndpoint[source]; exists {
		sourceToEndpoint[source] = ""
		sourceToParser[source] = nil
	}
}
