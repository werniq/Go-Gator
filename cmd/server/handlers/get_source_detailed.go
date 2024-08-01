package handlers

import (
	"github.com/gin-gonic/gin"
	"gogator/cmd/parsers"
	"net/http"
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

	c.JSON(http.StatusOK, gin.H{
		"sources": parsers.GetSourceDetailed(source),
	})
}
