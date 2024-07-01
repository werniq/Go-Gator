package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"newsaggr/cmd/parsers"
	"newsaggr/cmd/types"
)

// UpdateSource updates existent source with given parameters.
// If not-existent source is going to be updated - throws an error.
func UpdateSource(c *gin.Context) {
	var reqBody types.Source

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrFailedToDecode + err.Error(),
		})
		return
	}

	if reqBody.Name != "" {
		if reqBody.Endpoint != "" {
			parsers.UpdateSourceEndpoint(reqBody.Name, reqBody.Endpoint)
		}
		if reqBody.Format != "" {
			parsers.UpdateSourceFormat(reqBody.Name, reqBody.Format)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": MsgSourceUpdated,
	})
}
