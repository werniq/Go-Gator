package validator

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
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
	ErrDateFromAfter = "date from can not be after date end"

	// ErrFailedDateValidation is thrown when user submitted date in wrong format
	ErrFailedDateValidation = "error while validating date. correct format is YYYY-mm-dd - 2024-05-15"

	// ErrFailedSourceValidation is thrown when user submitted wrong source
	ErrFailedSourceValidation = "error while validating source "
)

// Handler interface defines the methods required for implementing a validation handler
type Handler interface {
	// SetNext sets the next handler in the chain

	SetNext(handler Handler)
	// Handle processes the request and performs validation
	Handle(c *gin.Context) error
}

// BaseHandler provides a base implementation for chaining handlers
type BaseHandler struct {
	next Handler
}

// SetNext sets the next handler in the chain
func (h *BaseHandler) SetNext(handler Handler) {
	h.next = handler
}

// HandleNext calls the next handler in the chain, if it exists
func (h *BaseHandler) HandleNext(c *gin.Context) error {
	if h.next != nil {
		return h.next.Handle(c)
	}
	return nil
}

// DateRangeHandler validates that the date-from parameter is not after the date-end parameter
type DateRangeHandler struct {
	BaseHandler
}

// Handle validates the date range and calls the next handler in the chain
func (h *DateRangeHandler) Handle(c *gin.Context) error {
	dateFrom := c.Query(DateFromFlag)
	dateEnd := c.Query(DateEndFlag)

	if dateEnd != "" && dateFrom != "" {
		if dateFrom > dateEnd {
			return errors.New(ErrDateFromAfter)
		}
	}

	return h.HandleNext(c)
}

// DateValidationHandler checks if the date parameters are in the correct format (YYYY-MM-DD)
type DateValidationHandler struct {
	BaseHandler
}

// Handle validates the date format and calls the next handler in the chain
func (h *DateValidationHandler) Handle(c *gin.Context) error {
	dateFrom := c.Query(DateFromFlag)
	dateEnd := c.Query(DateEndFlag)

	if err := ByDate(dateFrom); err != nil {
		return errors.New(ErrFailedDateValidation)
	}

	if err := ByDate(dateEnd); err != nil {
		return errors.New(ErrFailedDateValidation)
	}

	return h.HandleNext(c)
}

// SourceValidationHandler checks if the provided sources are within the supported list
type SourceValidationHandler struct {
	BaseHandler
}

// Handle validates the sources and calls the next handler in the chain
func (h *SourceValidationHandler) Handle(c *gin.Context) error {
	sources := c.Query(SourcesFlag)

	if err := BySources(sources); err != nil {
		return errors.New(ErrFailedSourceValidation)
	}

	return h.HandleNext(c)
}

// ByDate checks if the date string is in the correct format (YYYY-MM-DD)
func ByDate(dateStr string) error {
	if dateStr == "" {
		return nil
	}

	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format for %s, expected YYYY-MM-DD", dateStr)
	}

	return nil
}

// BySources checks if the provided sources are within the supported list
func BySources(sources string) error {
	if sources == "" {
		return nil
	}
	supportedSources := []string{"abc", "bbc", "nbc", "usatoday", "washingtontimes"}

	for _, source := range strings.Split(sources, ",") {
		if !contains(supportedSources, source) {
			return fmt.Errorf("unsupported source: %s. Supported sources are: %v", source, supportedSources)
		}
	}
	return nil
}

// contains checks if a slice contains a given string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}
