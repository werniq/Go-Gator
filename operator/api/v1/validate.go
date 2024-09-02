package v1

import (
	"errors"
	"net/url"
	"time"
)

const (
	// urlValidationError is a constant that represents the error message for invalid url
	urlValidationError = "url must contain http or https"
)

// validateFeeds function initializes a chain from all existing validation handlers
// and returns an error if any of the handlers fails
func validateFeeds(feed FeedSpec) error {
	urlValidationHandler := &urlValidate{url: feed.Link}

	return urlValidationHandler.Validate()
}

func validateHotNews(hotNewsSpec HotNewsSpec) error {
	dateValidationHandler := &dateValidate{dateStart: hotNewsSpec.DateStart, dateEnd: hotNewsSpec.DateEnd}

	return dateValidationHandler.Validate()
}

type handler interface {
	Validate() error
	SetNext(handler handler) handler
	HandleNext() error
}

type baseHandler struct {
	next handler
}

func (h *baseHandler) SetNext(handler handler) handler {
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

// Validate checks if the url is valid by using url.Parse
func (u *urlValidate) Validate() error {
	link, err := url.Parse(u.url)
	if err != nil {
		return err
	}
	if link.Scheme != "" && link.Host != "" {
		return u.HandleNext()
	}

	return errors.New(urlValidationError)
}

// dateValidate struct is used to check if the start date is before the end date
// and if the date format is correct
type dateValidate struct {
	baseHandler
	dateStart string
	dateEnd   string
}

func (d *dateValidate) Validate() error {
	if d.dateStart > d.dateEnd {
		return errors.New(errInvalidDateRange)
	}

	_, err := time.Parse(time.DateOnly, d.dateStart)
	if err != nil {
		return errors.New("invalid start date format: " + err.Error())
	}

	_, err = time.Parse(time.DateOnly, d.dateEnd)
	if err != nil {
		return errors.New("invalid end date format: " + err.Error())
	}

	return d.HandleNext()
}
