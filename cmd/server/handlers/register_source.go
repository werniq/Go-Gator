package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"newsaggr/cmd/parsers"
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
	var reqBody struct {
		Source   string `json:"source"`
		Endpoint string `json:"endpoint"`
		Format   string `json:"format"`
	}

	err := c.ShouldBindJSON(&reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error ": ErrFailedToDecode + err.Error(),
		})
		return
	}

	if sourceInArray(reqBody.Source) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrSourceExists,
		})
		return
	}

	parsers.AddNewSource(reqBody.Format, reqBody.Source, reqBody.Endpoint)

	c.JSON(http.StatusCreated, gin.H{
		"status": MsgSourceCreated,
	})
}
