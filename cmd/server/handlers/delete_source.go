package handlers

import (
	"github.com/gin-gonic/gin"
	"gogator/cmd/parsers"
	"log"
	"net/http"
)

const (
	// ErrDeleteSource is thrown whenever we encounter error while deleting new source (Admin API)
	ErrDeleteSource = "Failed to delete source: "

	// MsgSourceDeleted displays informational message after source was removed
	MsgSourceDeleted = "Source was successfully removed."
)

// DeleteSource handler deletes existing source from registered sources.
// If non-existent source is going to be deleted - throws an error.
func DeleteSource(c *gin.Context) {
	var reqBody struct {
		Source string `json:"name"`
	}

	err := c.ShouldBindJSON(&reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": ErrFailedToDecode + err.Error(),
		})
		return
	}

	if !sourceInArray(reqBody.Source) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrSourceNotFound,
		})
		return
	}

	err = parsers.DeleteSource(reqBody.Source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": ErrDeleteSource + err.Error(),
		})
		log.Println(ErrDeleteSource + err.Error())
		return
	}

	dates, err := parsers.GenerateDateRange(FirstFetchedFileDate, LastFetchedFileDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrDeleteSource + err.Error(),
		})
		log.Println(ErrDeleteSource + err.Error())
		return
	}

	err = parsers.DestroySource(reqBody.Source, dates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrDeleteSource + err.Error(),
		})
		log.Println(ErrDeleteSource + err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": MsgSourceDeleted,
	})
}
