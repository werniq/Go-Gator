package validator

import (
	"errors"
	"fmt"
	parsers "gogator/cmd/parsers"
	"strings"
	"time"
)

const (
	// ErrDateFromAfter is thrown when user provided date bigger than dateEnd
	ErrDateFromAfter = "date from can not be after date end"

	// ErrFailedDateValidation is thrown when user submitted date in wrong format
	ErrFailedDateValidation = "error while validating date. correct format is YYYY-mm-dd - 2024-05-15"

	// ErrFailedSourceValidation is thrown when user submitted wrong source
	ErrFailedSourceValidation = "error while validating source "
)

type Validator interface {
	Validate(keywords, dateFrom, dateEnd string) error
}

// ArgValidator struct is used to validate arguments which user inputs in get-news handler
type ArgValidator struct {
}

// Validate checks if all user-given arguments are correct
func (v *ArgValidator) Validate(sources, dateFrom, dateEnd string) error {
	dateFromValidator := &DateValidationHandler{
		date: dateFrom,
	}
	dateEndValidator := &DateValidationHandler{
		date: dateEnd,
	}
	dateRangeValidator := &DateRangeHandler{
		dateFrom: dateFrom,
		dateEnd:  dateEnd,
	}
	sourcesValidator := &SourceValidationHandler{
		sources: sources,
	}

	dateFromValidator.SetNext(dateEndValidator)
	dateEndValidator.SetNext(dateRangeValidator)
	dateRangeValidator.SetNext(sourcesValidator)

	if err := dateFromValidator.Handle(); err != nil {
		return err
	}

	return nil
}

// Handler interface defines the methods required for implementing a validation handler
type Handler interface {
	// SetNext sets the next handler in the chain
	SetNext(handler Handler)

	// Handle processes the request and performs validation
	Handle() error
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
func (h *BaseHandler) HandleNext() error {
	if h.next != nil {
		return h.next.Handle()
	}
	return nil
}

// DateRangeHandler validates that the date-from parameter is not after the date-end parameter
type DateRangeHandler struct {
	BaseHandler
	dateFrom string
	dateEnd  string
}

// Handle validates the date range and calls the next handler in the chain
func (h *DateRangeHandler) Handle() error {
	if err := ByDateRange(h.dateFrom, h.dateEnd); err != nil {
		return err
	}

	return h.HandleNext()
}

// DateValidationHandler checks if the date parameters are in the correct format (YYYY-MM-DD)
type DateValidationHandler struct {
	BaseHandler
	date string
}

// Handle validates the date format and calls the next handler in the chain
func (h *DateValidationHandler) Handle() error {
	if err := ByDate(h.date); err != nil {
		return errors.New(ErrFailedDateValidation)
	}

	return h.HandleNext()
}

// SourceValidationHandler checks if the provided sources are within the supported list
type SourceValidationHandler struct {
	BaseHandler
	sources string
}

// Handle validates the sources and calls the next handler in the chain
func (h *SourceValidationHandler) Handle() error {
	if err := BySources(h.sources); err != nil {
		return errors.New(ErrFailedSourceValidation + err.Error())
	}

	return h.HandleNext()
}

// ByDateRange verifies if date range is correct: date from should be before date end
func ByDateRange(dateFrom, dateEnd string) error {
	if dateFrom != "" && dateEnd != "" {
		if dateFrom > dateEnd {
			return errors.New(ErrDateFromAfter)
		}
	}

	return nil
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

	supportedSources := parsers.GetAllSources()

	for _, source := range strings.Split(sources, ",") {
		if _, exists := supportedSources[source]; !exists {
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
