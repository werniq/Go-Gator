package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeleteSource handler deletes existing source from registered sources.
// If non-existent source is going to be deleted - throws an error.
func DeleteSource(c *gin.Context) {
	var reqBody struct {
		Source string `json:"source"`
	}

	err := c.ShouldBindJSON(&reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error ": ErrFailedToDecode + err.Error(),
		})
		return
	}

	if !sourceInArray(reqBody.Source) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrSourceNotFound,
		})
		return
	}

	for i := 0; i <= len(Sources)-1; i++ {
		if Sources[i] == reqBody.Source {
			Sources = append(Sources[:i], Sources[i+1:]...)
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": MsgSourceDeleted,
	})
}
