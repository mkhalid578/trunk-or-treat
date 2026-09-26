package main

import (
	"car-api/internal/config"
	"car-api/internal/database"
	"car-api/internal/handlers"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer pool.Close()

	router := gin.Default()
	gin.SetMode(cfg.Mode)

	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.GET("/", func(c *gin.Context) {
		// If the client is 192.168.1.2, use the X-Forwarded-For
		// header to deduce the original client IP from the trust-
		// worthy parts of that header.
		// Otherwise, simply return the direct client IP
		fmt.Printf("ClientIP: %s\n", c.ClientIP())
	})
	router.POST("/auth/register", handlers.CreateUserHandler(pool))

	router.POST("/cars", handlers.CreateCarHandler(pool))
	router.GET("/cars", handlers.GetAllCarsHandler(pool))
	router.GET("/cars/:id", handlers.GetCarByIDHandler(pool))
	router.PUT("/cars/:id", handlers.UpdateCarVolumeHandler(pool))
	router.DELETE("/cars/:id", handlers.DeleteCarHandler(pool))

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
