package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetNews handler will be used in our server to retrieve news from files.
func GetNews(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
	})
}
