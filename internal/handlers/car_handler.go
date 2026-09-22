package handlers

import (
	"car-api/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

func GetAllCarsHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		cars, err := repository.GetAllCars(pool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, cars)
	}
}

func GetCarByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid car ID"})
			return
		}
		car, err := repository.GetCarByID(pool, id)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return

		}

		c.JSON(http.StatusOK, car)
	}
}

func UpdateCarVolumeHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid car ID"})
			return
		}

		var input struct {
			CargoFullVolumeCuFt   float64 `json:"cargo_volume_cu_ft"`
			CargoBehind2ndRowCuFt float64 `json:"cargo_vol_behind_2nd_row_cu_ft"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		car, err := repository.UpdateCarVolume(pool, id, input.CargoBehind2ndRowCuFt, input.CargoFullVolumeCuFt)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, car)
	}
}
