package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"newsaggr/cmd/filters"
	"newsaggr/cmd/parsers"
	"newsaggr/cmd/types"
	"newsaggr/cmd/validator"
	"time"
)

var (
	// LastFetchedFileDate will be used for iterating over files with news
	LastFetchedFileDate = time.Now().Format(time.DateOnly)
)

const (
	// KeywordFlag will be used to get the keywords (or empty string) from URL parameter
	KeywordFlag = "keywords"

	// DateFromFlag will be used to get the date-from (or empty string) from URL parameter
	DateFromFlag = "date-from"

	// DateEndFlag will be used to get the date-end (or empty string) from URL parameter
	DateEndFlag = "date-end"

	// SourcesFlag will be used to get the sources (or empty string) from URL parameter
	SourcesFlag = "sources"

	// ErrFailedParsing is thrown when program fails to parse sources
	ErrFailedParsing = "error while parsing sources: "

	// FirstFetchedFileDate identifies first file which contains news
	FirstFetchedFileDate = "2024-07-05"
)

// GetNews handler will be used in our server to retrieve news from prepared files
func GetNews(c *gin.Context) {
	keywords := c.Query(KeywordFlag)
	sources := c.Query(SourcesFlag)
	dateFrom := c.Query(DateFromFlag)
	dateEnd := c.Query(DateEndFlag)

	dateRangeHandler := &validator.DateRangeHandler{}
	dateValidationHandler := &validator.DateValidationHandler{}
	sourceValidationHandler := &validator.SourceValidationHandler{}

	dateRangeHandler.SetNext(dateValidationHandler)
	dateValidationHandler.SetNext(sourceValidationHandler)

	// Start the chain
	if err := dateRangeHandler.Handle(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	params := types.NewFilteringParams(keywords, dateFrom, dateEnd, sources)

	var news []types.News
	if dateFrom == "" {
		dateFrom = FirstFetchedFileDate
	}
	if dateEnd == "" {
		dateEnd = LastFetchedFileDate
	}

	news, err := parsers.FromFiles(dateFrom, dateEnd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	news = filters.Apply(news, params)

	c.JSON(http.StatusOK, gin.H{
		"totalAmount": len(news),
		"news":        news,
	})
}
