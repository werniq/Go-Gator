package handlers

import (
	"github.com/gin-gonic/gin"
	"gogator/cmd/parsers"
	"gogator/cmd/types"
	"log"
	"net/http"
)

const (

	// MsgSourceUpdated returns a successful message after changing source information
	MsgSourceUpdated = "Source was successfully updated"

	// ErrSourceNotFound displays when we try to delete not-existent source
	ErrSourceNotFound = "Source is not found in available sources. Please, check the name and try again."

	// ErrUpdateSource is thrown whenever we encounter error while updating new source (Admin API)
	ErrUpdateSource = "Failed to update source: "

	ErrNoSourceName = "No source name detected. Please, provide source name."
)

// UpdateSource updates existent source with given parameters.
// If not-existent source is going to be updated - throws an error.
func UpdateSource(c *gin.Context) {
	var reqBody types.Source
	var err error

	if err = c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrFailedToDecode + err.Error(),
		})
		return
	}

	if reqBody.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrNoSourceName,
		})
		return
	}

	if !sourceInArray(reqBody.Name) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrSourceNotFound,
		})
		return
	}

	if reqBody.Endpoint != "" {
		err = parsers.UpdateSourceEndpoint(reqBody.Name, reqBody.Endpoint)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": ErrUpdateSource + err.Error(),
			})
			log.Println(ErrUpdateSource + err.Error())
			return
		}
	}

	if reqBody.Format != "" {
		err = parsers.UpdateSourceFormat(reqBody.Name, reqBody.Format)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": ErrUpdateSource + err.Error(),
			})
			log.Println(ErrUpdateSource + err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": MsgSourceUpdated,
	})
}
