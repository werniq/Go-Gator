package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"newsaggr/cmd/parsers"
)

// UpdateSource updates existent source with given parameters.
// If not-existent source is going to be updated - throws an error.
func UpdateSource(c *gin.Context) {
	var reqBody struct {
		Source   string `json:"source"`
		Endpoint string `json:"endpoint"`
		Format   string `json:"format"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrFailedToDecode + err.Error(),
		})
		return
	}

	if reqBody.Source != "" {
		if reqBody.Endpoint != "" {
			parsers.UpdateSourceEndpoint(reqBody.Source, reqBody.Endpoint)
		}
		if reqBody.Format != "" {
			parsers.UpdateSourceFormat(reqBody.Source, reqBody.Format)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": MsgSourceUpdated,
	})
}
