package main

import (
	"log"
	"os"

	"github.com/Anom-a/Orchestrix/internal/config"
	"github.com/Anom-a/Orchestrix/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../../.env")

	port := os.Getenv("PORT")
	configFile := config.Load(port)

	r := gin.Default()
	routes.RegisterRoutes(r)

	log.Printf("Starting server on %s", configFile.Port)
	if err := r.Run(configFile.Port); err != nil {
		log.Fatal(err.Error())
	}
}
