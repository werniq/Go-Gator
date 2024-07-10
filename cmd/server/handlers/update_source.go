package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"newsaggr/cmd/parsers"
	"newsaggr/cmd/types"
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

	if sourceInArray(reqBody.Name) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrSourceExists,
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
