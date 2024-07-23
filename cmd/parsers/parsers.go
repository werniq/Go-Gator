package parsers

import (
	"gogator/cmd/types"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Parser interface will be used to implement parsers
//
// It provides a Parse method, which is used to open json, xml or html file
// retrieve all data from it, and parse into an array of articles.
// Returns an error, if file not exists, or error while decoding data.
type Parser interface {
	Parse() ([]types.News, error)
}

type getFileData func(f string) ([]byte, error)

var (
	StoragePath string

	openFile getFileData = extractFileData

	// g and f are Parsing and instruction factories.
	// These are custom types which will be used for parsers initialization for
	// different data formats, and applying filters to retrieved data
	g ParsingFactory

	// sourceToEndpoint maps source names (as strings) to their corresponding filenames
	sourceToEndpoint = map[string]string{
		WashingtonTimes: "https://www.washingtontimes.com/rss/headlines/news/world",
		ABC:             "https://abcnews.go.com/abcnews/internationalheadlines",
		BBC:             "https://feeds.bbci.co.uk/news/rss.xml",
		UsaToday:        "https://usatoday.com",
	}

	// sourceToParser maps source names to their corresponding parser instances
	// The parser are created using ParsingFactory
	sourceToParser = map[string]Parser{
		UsaToday:        g.HtmlParser(UsaToday),
		ABC:             g.XmlParser(ABC),
		BBC:             g.XmlParser(BBC),
		WashingtonTimes: g.XmlParser(WashingtonTimes),
	}
)

const (
	// Define Sources
	NbcNews         = "nbc"
	UsaToday        = "usatoday"
	ABC             = "abc"
	BBC             = "bbc"
	WashingtonTimes = "washingtontimes"

	AllSources = ""

	CmdDir = "cmd"

	ParsersDir = "parsers"

	DataDir = "data"

	sourcesFile = "sources" + JsonExtension

	JsonExtension = ".json"
)

// ParseBySource returns all news in particular source. If source is equal to "all", news will be
// retrieved from all sources
func ParseBySource(source string) ([]types.News, error) {
	var (
		news       []types.News
		wg         sync.WaitGroup
		mu         sync.Mutex
		errChannel = make(chan error, 1)
	)

	collectNews := func(p Parser) {
		defer wg.Done()
		n, err := p.Parse()
		if err != nil {
			select {
			case errChannel <- err:
			default:
			}
			return
		}
		mu.Lock()
		news = append(news, n...)
		mu.Unlock()
	}

	if source == "" {
		for _, p := range sourceToParser {
			wg.Add(1)
			go collectNews(p)
		}
	} else {
		sources := strings.Split(source, ",")
		for _, s := range sources {
			if p, exists := sourceToParser[s]; exists {
				wg.Add(1)
				go collectNews(p)
			}
		}
	}

	wg.Wait()
	close(errChannel)

	if err, ok := <-errChannel; ok {
		return nil, err
	}

	return news, nil
}

// extractFileData reads data from file $filename and returns its content
func extractFileData(filename string) ([]byte, error) {
	cwdPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	f := filepath.Join(cwdPath, StoragePath, filename)

	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	return data, nil
}
