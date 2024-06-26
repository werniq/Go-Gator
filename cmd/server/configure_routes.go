package server

import (
	"github.com/gin-gonic/gin"
	"newsaggr/cmd/server/handlers"
)

// setupRoutes attaches routes to *gin.Engine
func setupRoutes(r *gin.Engine) {
	r.GET("/news", handlers.GetNews)

	admin := r.Group("/admin")

	sources := admin.Group("/sources")
	sources.GET("/", handlers.GetSources)
	sources.POST("/", handlers.RegisterSource)
	sources.PUT("/", handlers.UpdateSource)
	sources.DELETE("/", handlers.DeleteSource)
}
