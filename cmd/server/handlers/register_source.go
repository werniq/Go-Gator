package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"newsaggr/cmd/parsers"
	"newsaggr/cmd/types"
)

// sourceInArray checks if sources is already in array
func sourceInArray(source string) bool {
	if _, exists := parsers.GetAllSources()[source]; exists {
		return true
	}
	return false
}

// RegisterSource handler will be used in order to create new source from where
// we can parse news
func RegisterSource(c *gin.Context) {
	var reqBody types.Source

	err := c.ShouldBindJSON(&reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error ": ErrFailedToDecode + err.Error(),
		})
		return
	}

	if sourceInArray(reqBody.Name) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrSourceExists,
		})
		return
	}

	err = parsers.AddNewSource(reqBody.Format, reqBody.Name, reqBody.Endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": ErrAddSource + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": MsgSourceCreated,
	})
}
