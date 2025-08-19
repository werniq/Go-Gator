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
	Parse() ([]types.Article, error)
}

const (
	// UsaToday represents the identifier for USA Today source
	UsaToday = "usatoday"

	// ABC represents the identifier for ABC source
	ABC = "abc"

	// BBC represents the identifier for BBC source
	BBC = "bbc"

	// WashingtonTimes represents the identifier for Washington Times source
	WashingtonTimes = "washingtontimes"

	// AllSources is an empty string used to indicate all sources
	AllSources = ""

	// CmdDir is a directory where ParsersDir is located
	CmdDir = "cmd"

	// ParsersDir is a directory where DataDir is located
	ParsersDir = "parsers"

	// DataDir is the directory name where data files are stored
	DataDir = "data"

	// sourcesFile is the filename for the sources JSON file
	sourcesFile = "sources" + JsonExtension

	// JsonExtension is the file extension used for JSON files
	JsonExtension = ".json"
)

// ParseBySource retrieves all news from a particular source.
//
// If the source parameter is equal to "all", news will be retrieved from all sources specified in sourceToParser.
//
// The function returns a slice of news items and an error if any occurred during the parsing process.
func ParseBySource(source string) ([]types.Article, error) {
	var (
		news       []types.Article
		wg         sync.WaitGroup
		mu         sync.Mutex
		errChannel = make(chan error, 1)
	)

	if source == "" {
		for _, p := range sourceToParser {
			wg.Add(1)
			go fetchNews(p, &news, &wg, &mu, errChannel)
		}
	} else {
		sources := strings.Split(source, ",")
		for _, sourceName := range sources {
			if p, exists := sourceToParser[sourceName]; exists {
				wg.Add(1)
				go fetchNews(p, &news, &wg, &mu, errChannel)
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

// fetchNews is a helper function to parse news from a given parser.
//
// # It updates the news slice in a concurrency-safe manner and sends any errors to errChannel
//
// We use pointers to all variables from function ParseBySource and FromFiles.
// It will cause a panic if we will call wg.Done() without passing a pointer:
// / each goroutine would receive its own copy of the WaitGroup, which leads to incorrect synchronization:
// / because the Add, Done, and Wait calls would affect separate WaitGroup instances,
// / and most likely causing the Wait() function to never return or behave unpredictably.
func fetchNews(p Parser, news *[]types.Article, wg *sync.WaitGroup, mu *sync.Mutex, errChannel chan<- error) {
	defer wg.Done()

	parsedNews, err := p.Parse()
	if err != nil {
		errChannel <- err
		return
	}

	mu.Lock()
	*news = append(*news, parsedNews...)
	mu.Unlock()
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
