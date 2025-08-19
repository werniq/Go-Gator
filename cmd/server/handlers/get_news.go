package handlers

import (
	"github.com/gin-gonic/gin"
	"gogator/cmd/filters"
	"gogator/cmd/parsers"
	"gogator/cmd/types"
	"gogator/cmd/validator"
	"log"
	"net/http"
	"time"
)

var (
	// LastFetchedFileDate will be used for iterating over files with news
	//
	// It is assigned date of launching the server and will be often updated
	LastFetchedFileDate = time.Now().Format(time.DateOnly)

	// FirstFetchedFileDate identifies first file which contains news
	//
	// It is assigned date of launching the server, and won't be updated.
	FirstFetchedFileDate = time.Now().Format(time.DateOnly)
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

	//
	ErrValidatingParams = "Error validating parameters: "
)

// GetNews handler will be used in our server to retrieve news from prepared files
func GetNews(c *gin.Context) {
	keywords := c.Query(KeywordFlag)
	sources := c.Query(SourcesFlag)
	dateFrom := c.Query(DateFromFlag)
	dateEnd := c.Query(DateEndFlag)

	v := &validator.ArgValidator{}
	err := v.Validate(sources, dateFrom, dateEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrValidatingParams + err.Error(),
		})
		log.Println(ErrValidatingParams + err.Error())
		return
	}

	params := types.NewFilteringParams(keywords, dateFrom, dateEnd, sources)

	var news []types.Article
	if dateFrom == "" {
		dateFrom = FirstFetchedFileDate
	}
	if dateEnd == "" {
		dateEnd = LastFetchedFileDate
	}

	news, err = parsers.FromFiles(dateFrom, dateEnd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": ErrFailedParsing + err.Error(),
		})
		log.Println(ErrFailedParsing, err.Error())
		return
	}

	news = filters.Apply(news, params)

	c.JSON(http.StatusOK, gin.H{
		"totalAmount": len(news),
		"news":        news,
	})
}
