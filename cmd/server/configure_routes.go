package server

import (
	"github.com/gin-gonic/gin"
	"gogator/cmd/server/handlers"
)

// setupRoutes attaches routes to *gin.Engine
func setupRoutes(r *gin.Engine) {
	r.GET("/news", handlers.GetNews)

	r.GET("/admin/sources", handlers.GetSources)
	r.GET("/admin/sources/:source", handlers.GetSourceDetailed)
	r.POST("/admin/sources", handlers.RegisterSource)
	r.PUT("/admin/sources", handlers.UpdateSource)
	r.DELETE("/admin/sources", handlers.DeleteSource)
}
