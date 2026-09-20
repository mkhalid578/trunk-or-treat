package handlers

import (
	"car-api/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateCarInput struct {
	Make       string `json:"make" binding:"required"`
	Model      string `json:"model" binding:"required"`
	Year       int    `json:"model_year" binding:"required"`
	Trim       string `json:"trim" binding:"required"`
	BodyStyle  string `json:"body_style" binding:"required"`
	Powertrain string `json:"powertrain" binding:"required"`
}

func CreateCarHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateCarInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		car, err := repository.AddCar(
			pool,
			input.Model,
			input.Make,
			input.Year,
			input.Trim,
			input.BodyStyle,
			input.Powertrain,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, car)
	}
}
