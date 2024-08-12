package v1

import (
	"fmt"
	"strings"
)

const (
	// urlValidationError is a constant that represents the error message for invalid url
	urlValidationError = "url must contain http or https"
)

// Validate function initializes a chain from all existing validation handlers
// and returns an error if any of the handlers fails
func Validate(feed FeedSpec) error {
	urlValidationHandler := &urlValidate{url: feed.Link}

	return urlValidationHandler.Validate()
}

type Handler interface {
	Validate() error
	SetNext(handler Handler) Handler
	HandleNext() error
}

type baseHandler struct {
	next Handler
}

func (h *baseHandler) SetNext(handler Handler) Handler {
	h.next = handler
	return handler
}

func (h *baseHandler) HandleNext() error {
	if h.next != nil {
		return h.next.Validate()
	}
	return nil
}

// urlValidate struct is used to check if request url is constructed properly
type urlValidate struct {
	baseHandler
	url string
}

// Validate checks if the url contains http or https
func (u *urlValidate) Validate() error {
	// check by regular expression if the url contains http or https
	if !strings.Contains(u.url, "http") {
		return fmt.Errorf(urlValidationError)
	}
	return u.HandleNext()
}
