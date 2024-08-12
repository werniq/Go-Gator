package v1

import (
	"fmt"
	"strings"
)

const (
	// lengthValidationError is a constant that represents the error message for invalid length of keyword
	lengthValidationError = "length of keyword is invalid: "

	// urlValidationError is a constant that represents the error message for invalid url
	urlValidationError = "url must contain http or https"
)

// Validate function initializes a chain from all existing validation handlers
// and returns an error if any of the handlers fails
func Validate(feed FeedSpec) error {
	urlValidationHandler := &UrlValidate{url: feed.Link}

	return urlValidationHandler.Validate()
}

type Handler interface {
	Validate() error
	SetNext(handler Handler) Handler
	HandleNext() error
}

type BaseHandler struct {
	next Handler
}

func (h *BaseHandler) SetNext(handler Handler) Handler {
	h.next = handler
	return handler
}

func (h *BaseHandler) HandleNext() error {
	if h.next != nil {
		return h.next.Validate()
	}
	return nil
}

// LengthValidate struct is used to check if the length of keyword is within the required range
type LengthValidate struct {
	BaseHandler
	requiredMinLength int
	requiredMaxLength int
	keyword           string
}

// Validate func in LengthValidate struct is used to check if the length of keyword is within the required range
func (l *LengthValidate) Validate() error {
	if len(l.keyword) < l.requiredMinLength ||
		len(l.keyword) > l.requiredMaxLength {
		return fmt.Errorf("%s %s(%d)", lengthValidationError, l.keyword, len(l.keyword))
	}
	return l.HandleNext()
}

// UrlValidate struct is used to check if request url is constructed properly
type UrlValidate struct {
	BaseHandler
	url string
}

// Validate checks if the url contains http or https
func (u *UrlValidate) Validate() error {
	if !strings.Contains(u.url, "http") {
		return fmt.Errorf(urlValidationError)
	}
	return u.HandleNext()
}
