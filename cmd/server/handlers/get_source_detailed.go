package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"newsaggr/cmd/parsers"
)

// GetSourceDetailed returns detailed information about source
func GetSourceDetailed(c *gin.Context) {
	source := c.Param("source")

	if source == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrNoSourceName,
		})
		return
	}

	if !sourceInArray(source) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrSourceNotFound,
		})
		return
	}

	c.JSON(200, gin.H{
		"sources": parsers.GetSourceDetailed(source),
	})
}
