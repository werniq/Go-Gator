package v1

import (
	"fmt"
	"strings"
)

const (
	lengthValidationError = "length of keyword is invalid: "

	nameValidationError = "name cannot be empty"

	urlValidationError = "url must contain http or https"

	requiredMinLength = 1

	requiredMaxLength = 20
)

// Validate function initializes a chain from all existing validation handlers
// and returns an error if any of the handlers fails
func Validate(feed *Feed) error {
	lengthValidationHandler := &LengthValidate{
		keyword:           feed.Spec.Name,
		requiredMaxLength: requiredMaxLength,
		requiredMinLength: requiredMinLength,
	}

	urlValidationHandler := &UrlValidate{url: feed.Spec.Link}

	lengthValidationHandler.
		SetNext(urlValidationHandler)

	return lengthValidationHandler.Validate()
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
