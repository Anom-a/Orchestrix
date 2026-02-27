package routes

import (
	"github.com/Anom-a/Orchestrix/internal/handlers"
	"github.com/Anom-a/Orchestrix/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine){
	api := r.Group("/api")
	api.GET("/health", handlers.HealthCheck)
	auth := r.Group("/auth")
	auth.POST("/login", handlers.Login)
	auth.POST("/register", handlers.Register)
	documents := api.Group("/documents")
	documents.Use(middleware.JWTAuth())
	documents.POST("/upload", handlers.UploadDocument)
	documents.GET("/", handlers.ListDocuments)
}