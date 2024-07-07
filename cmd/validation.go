package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
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
)

type UnsupportedFlagError struct {
	Err error
}

func (u *UnsupportedFlagError) Error() string {
	return fmt.Sprintf(u.Err.Error())
}

// Handler interface defines the methods required for implementing a validation handler
type Handler interface {
	// Handle verifies if given string is following rules of this handler
	Handle(cmd *cobra.Command) error

	// SetNext sets the next handler in the chain
	SetNext(handler Handler)
}

type BaseHandler struct {
	next Handler
}

// HandleNext passes the command to the next handler in the chain if it exists
func (b *BaseHandler) HandleNext(cmd *cobra.Command) error {
	if b.next != nil {
		return b.next.Handle(cmd)
	}

	return nil
}

// SetNext sets the next handler in the chain
func (b *BaseHandler) SetNext(handler Handler) {
	b.next = handler
}

type DateValidationHandler struct {
	BaseHandler
}

// Handle validates the date flags and checks their logical consistency
func (h *DateValidationHandler) Handle(cmd *cobra.Command) error {
	dateFrom, err := cmd.Flags().GetString(DateFromFlag)
	err = checkFlagErr(err)
	if err != nil {
		return err
	}

	err = validateDate(dateFrom)
	if err != nil {
		return err
	}

	dateEnd, err := cmd.Flags().GetString(DateEndFlag)
	if err != nil {
		return err
	}

	err = validateDate(dateEnd)
	if err != nil {
		return err
	}

	// Ensure that dateFrom is not after dateEnd
	if dateEnd != "" && dateFrom != "" {
		if dateFrom > dateEnd || dateEnd > dateFrom {
			log.Fatalln("Date from can not be after date end.")
		}
	}

	return h.HandleNext(cmd)
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
}

// Handle validates the sources flag and checks if the provided sources are within the supported list
func (h *SourcesValidationHandler) Handle(cmd *cobra.Command) error {
	sources, err := cmd.Flags().GetString(SourcesFlag)
	err = checkFlagErr(err)
	if err != nil {
		return err
	}

	err = validateSources(sources)
	if err != nil {
		return err
	}

	return h.HandleNext(cmd)
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

// checkFlagErr enhances flag-related error messages with more user-friendly versions
func checkFlagErr(err error) error {
	if err != nil {
		for substr, msg := range errorMessages {
			if strings.Contains(err.Error(), substr) {
				return &UnsupportedFlagError{Err: errors.New(msg + err.Error())}
			}
		}

		return errors.New("error parsing flags: " + err.Error())
	}

	return nil
}
