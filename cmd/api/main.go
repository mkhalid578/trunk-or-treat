package main

import (
	"car-api/internal/config"
	"car-api/internal/database"
	"car-api/internal/handlers"
	"car-api/internal/middleware"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
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
	router.Use(cors.Default())

	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.GET("/", func(c *gin.Context) {
		if err := pool.Ping(c.Request.Context()); err != nil {
			log.Printf("ping failed: %v", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.POST("/auth/register", handlers.CreateUserHandler(pool))
	router.POST("/auth/login", handlers.LoginHandler(pool, cfg))

	router.POST("/cars", middleware.AuthMiddleware(cfg), handlers.CreateCarHandler(pool))
	router.GET("/cars", handlers.GetAllCarsHandler(pool))
	router.GET("/cars/:id", middleware.AuthMiddleware(cfg), handlers.GetCarByIDHandler(pool))
	router.PUT("/cars/:id", middleware.AuthMiddleware(cfg), handlers.UpdateCarVolumeHandler(pool))
	router.DELETE("/cars/:id", middleware.AuthMiddleware(cfg), handlers.DeleteCarHandler(pool))

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
