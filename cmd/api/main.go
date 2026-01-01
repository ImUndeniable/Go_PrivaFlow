package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	//1. Initialize Router
	r := gin.Default()

	//2. Health Check Endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "active",
			"project": "PrivaFlow",
			"message": "Locked IN!!",
		})
	})

	// 3. Start Server
	log.Println("🚀 PrivaFlow Orchestrator running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
