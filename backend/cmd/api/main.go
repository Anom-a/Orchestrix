package main

import (
	"log"
	"os"
	"time"

	"github.com/Anom-a/Orchestrix/internal/config"
	"github.com/Anom-a/Orchestrix/internal/database"
	"github.com/Anom-a/Orchestrix/internal/models"
	"github.com/Anom-a/Orchestrix/internal/routes"
	"github.com/Anom-a/Orchestrix/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env", "../.env", "../../.env", "backend/.env")
	port := os.Getenv("PORT")
	configFile := config.Load(port)
	r := gin.Default()
	database.Connect()
	database.DB.AutoMigrate(&models.User{}, models.Document{})

	const staleProcessingTimeout = 20 * time.Minute
	recoveredCount, recoveryErr := services.RecoverStaleProcessingDocuments(staleProcessingTimeout)
	if recoveryErr != nil {
		log.Printf("Stale document recovery failed: %v", recoveryErr)
	} else if recoveredCount > 0 {
		log.Printf("Recovered %d stale processing documents", recoveredCount)
	}

	routes.RegisterRoutes(r)

	log.Printf("Starting server on %s", configFile.Port)
	if err := r.Run(configFile.Port); err != nil {
		log.Fatal(err.Error())
	}
}
