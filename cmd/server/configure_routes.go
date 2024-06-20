package server

import (
	"github.com/gin-gonic/gin"
	"newsaggr/cmd/server/handlers"
)

// setupRoutes attaches routes to *gin.Engine
func setupRoutes(r *gin.Engine) {
	r.GET("/", handlers.GetNews)

	admin := r.Group("/admin")
	admin.GET("/source", handlers.GetSources)
	admin.POST("/source", handlers.RegisterSource)
	admin.PUT("/source", handlers.UpdateSource)
	admin.DELETE("/source", handlers.DeleteSource)
}
