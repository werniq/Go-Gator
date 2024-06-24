package server

import (
	"github.com/gin-gonic/gin"
	"newsaggr/cmd/server/handlers"
)

// setupRoutes attaches routes to *gin.Engine
func setupRoutes(r *gin.Engine) {
	r.GET("/news", handlers.GetNews)

	admin := r.Group("/admin")

	source := admin.Group("/source")
	source.GET("/", handlers.GetSources)
	source.POST("/", handlers.RegisterSource)
	source.PUT("/", handlers.UpdateSource)
	source.DELETE("/", handlers.DeleteSource)
}
