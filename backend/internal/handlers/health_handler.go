package handlers

import "github.com/gin-gonic/gin"



func HealthCheck(c *gin.Context){
	c.IndentedJSON(200, gin.H{
		"status": "ok",
		"service": "orchestrix-backend",
	})
}