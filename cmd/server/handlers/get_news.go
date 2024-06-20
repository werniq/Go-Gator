package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"newsaggr/cmd/parsers"
	"newsaggr/cmd/types"
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

	// ErrDateFromAfter is thrown when user provided DateFrom bigger than DateEnd
	ErrDateFromAfter = "Date from can not be after date end"

	// ErrFailedParsing is thrown when program fails to parse sources
	ErrFailedParsing = "error while parsing sources: "
)

// GetNews handler will be used in our server to retrieve news from files.
func GetNews(c *gin.Context) {
	keywords := c.Query(KeywordFlag)
	sources := c.Query(SourcesFlag)
	dateFrom := c.Query(DateFromFlag)
	dateEnd := c.Query(DateEndFlag)

	if dateEnd != "" && dateFrom != "" {
		if dateFrom > dateEnd {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": ErrDateFromAfter,
			})
		}
	}

	params := types.NewFilteringParams(keywords, dateFrom, dateEnd)

	news, err := parsers.ParseBySource(sources)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": ErrFailedParsing + err.Error(),
		})
		return
	}

	news = parsers.ApplyFilters(news, params)

	c.JSON(http.StatusOK, gin.H{
		"totalLength": len(news),
		"news":        news,
	})
}
