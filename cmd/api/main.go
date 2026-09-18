package main

import (
	"car-api/internal/config"
	"car-api/internal/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// schedule a pool clean up
	// schedules close for the connection pool

	defer pool.Close()

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "pong",
			"database": "connected",
		})
	})

	router.Run(":" + cfg.Port)

}
