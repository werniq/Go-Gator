package handlers

import (
	"github.com/gin-gonic/gin"
	"newsaggr/cmd/parsers"
)

var (
	// MsgSourceCreated displays informational message after source was created
	MsgSourceCreated = "Source was successfully registered."

	// MsgSourceDeleted displays informational message after source was removed
	MsgSourceDeleted = "Source was successfully removed."

	// MsgSourceUpdated returns a successful message after changing source information
	MsgSourceUpdated = "Source was successfully updated"

	// ErrFailedToDecode displays when server failed to decode request body into struct
	ErrFailedToDecode = "Error while decoding request body: "

	// ErrSourceExists is throws when user tries to register already registered source
	ErrSourceExists = "This source is already registered. "

	// ErrSourceNotFound displays when we try to delete not-existent source
	ErrSourceNotFound = "Source is not found in available sources. Please, check the name and try again."
)

// GetSources returns all currently available news sources
func GetSources(c *gin.Context) {
	c.JSON(200, gin.H{
		"sources": parsers.GetAllSources(),
	})
}
