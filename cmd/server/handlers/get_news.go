package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"newsaggr/cmd/parsers"
	"newsaggr/cmd/types"
	"newsaggr/cmd/validator"
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

	// ErrFailedDateValidation is thrown when user submitted date in wrong format
	ErrFailedDateValidation = "error while validating date. correct format is YYYY-mm-dd - 2024-05-15"

	// ErrFailedSourceValidation is thrown when user submitted wrong source
	ErrFailedSourceValidation = "error while validating sources."
)

// GetNews handler will be used in our server to retrieve news from prepared files
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
			return
		}
	}

	err := validator.ByDate(dateFrom)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrFailedDateValidation + err.Error(),
		})
		return
	}

	err = validator.ByDate(dateEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrFailedDateValidation + err.Error(),
		})
		return
	}

	err = validator.BySources(sources)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrFailedSourceValidation + err.Error(),
		})
		return
	}

	params := types.NewFilteringParams(keywords, dateFrom, dateEnd)

	var news []types.News
	if dateFrom == "" {
		dateFrom = parsers.FirstFetchedFileDate
	}
	if dateEnd == "" {
		dateEnd = parsers.LastFetchedFileDate
	}

	news, err = parsers.FromFiles(dateFrom, dateEnd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	news = parsers.ApplyFilters(news, params)

	c.JSON(http.StatusOK, gin.H{
		"totalAmount": len(news),
		"news":        news,
	})
}
