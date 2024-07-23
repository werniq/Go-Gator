package handlers

import (
	"github.com/gin-gonic/gin"
	"gogator/cmd/parsers"
	"net/http"
)

// GetSources returns all currently available news sources
func GetSources(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"sources": parsers.GetAllSources(),
	})
}
