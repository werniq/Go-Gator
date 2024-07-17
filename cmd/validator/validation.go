package validator

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// errorMessages maps specific flag access errors to user-friendly error messages
	errorMessages = map[string]string{
		"flag accessed but not defined": "Unsupported flag: ",
	}

	// supportedSources defines the list of valid source values.
	supportedSources = []string{"abc", "bbc", "nbc", "usatoday", "washingtontimes", "all"}

	// ErrWrongSource is the error message for invalid source input.
	ErrWrongSource = fmt.Sprintf("Supported sources are: %v. Please, check spelling of inputed source: ", supportedSources)

	ErrInvalidDateRange = "invalid date range: make sure that date from is before date end"
)

// Validate checks if all mentioned flags are in correct format and
func Validate(dateFrom, dateEnd, sources string) error {
	sourcesValidationHandler := &SourcesValidationHandler{
		sources: sources,
	}
	dateFromValidationHandler := &DateValidationHandler{
		Date: dateFrom,
	}
	dateEndValidationHandler := &DateValidationHandler{
		Date: dateEnd,
	}
	dateRangeValidationHandler := &DateRangeValidatorHandler{
		DateFrom: dateFrom,
		DateEnd:  dateEnd,
	}

	sourcesValidationHandler.SetNext(dateFromValidationHandler)
	dateFromValidationHandler.SetNext(dateEndValidationHandler)
	dateEndValidationHandler.SetNext(dateRangeValidationHandler)

	err := sourcesValidationHandler.Handle()
	if err != nil {
		return err
	}

	return nil
}

type UnsupportedFlagError struct {
	Err error
}

func (u *UnsupportedFlagError) Error() string {
	return fmt.Sprintf(u.Err.Error())
}

// Handler interface defines the methods required for implementing a validation handler
type Handler interface {
	// Handle verifies if given string is following rules of this handler
	Handle() error

	// SetNext sets the next handler in the chain
	SetNext(handler Handler)
}

// BaseHandler provides functionality to set and handle next handler
type BaseHandler struct {
	next Handler
}

// HandleNext passes the command to the next handler in the chain if it exists
func (b *BaseHandler) HandleNext() error {
	if b.next != nil {
		return b.next.Handle()
	}

	return nil
}

// SetNext sets the next handler in the chain
func (b *BaseHandler) SetNext(handler Handler) {
	b.next = handler
}

// DateValidationHandler is used to validate dates format: it has to be YYYY-MM-DD
type DateValidationHandler struct {
	BaseHandler
	Date string
}

// Handle validates the date flags and checks their logical consistency
func (h *DateValidationHandler) Handle() error {
	err := validateDate(h.Date)
	if err != nil {
		return err
	}

	return h.HandleNext()
}

// DateRangeValidatorHandler
// after validating dates by format, we will check if date range is correct: dateFrom should be before dateEnd
type DateRangeValidatorHandler struct {
	BaseHandler
	DateFrom string
	DateEnd  string
}

func (h *DateRangeValidatorHandler) Handle() error {
	if h.DateFrom == "" || h.DateEnd == "" {
		return h.HandleNext()
	}

	if h.DateFrom > h.DateEnd {
		return errors.New(ErrInvalidDateRange)
	}
	// if dateEnd > dateFrom it will be just , so nil will be returned

	return h.HandleNext()
}

// CheckFlagErr enhances flag-related error messages with more user-friendly versions
func CheckFlagErr(err error) error {
	if err != nil {
		if strings.Contains(err.Error(), "flag accessed but not defined") {
			return &UnsupportedFlagError{Err: errors.New("Unsupported flag: " + err.Error())}
		}

		return errors.New("error parsing flags: " + err.Error())
	}

	return nil
}

// validateDate checks if the date string is in the correct format YYYY-MM-DD
func validateDate(dateStr string) error {
	if dateStr == "" {
		return nil
	}

	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format for %s, expected YYYY-MM-DD", dateStr)
	}

	return nil
}

type SourcesValidationHandler struct {
	BaseHandler
	sources string
}

// Handle validates the sources flag and checks if the provided sources are within the supported list
func (h *SourcesValidationHandler) Handle() error {
	err := validateSources(h.sources)
	if err != nil {
		return err
	}

	return h.HandleNext()
}

// validateSources checks if the provided sources are within the supported list
func validateSources(sources string) error {
	if sources == "" {
		return nil
	}

	for _, source := range strings.Split(sources, ",") {
		if !contains(supportedSources, source) {
			return errors.New(ErrWrongSource + source)
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
