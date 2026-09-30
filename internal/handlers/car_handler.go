package handlers

import (
	"car-api/internal/repository"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

type FilterRequest struct {
	Year       *int    `form:"year"`
	Make       *string `form:"make"`
	Model      *string `form:"model"`
	Trim       *string `form:"trim"`
	BodyStyle  *string `form:"body_style"`
	Powertrain *string `form:"powertrain"`
}

type UpdateCarVolume struct {
	CargoFullVolumeCuFt   *float64 `json:"cargo_volume_cu_ft"`
	CargoBehind2ndRowCuFt *float64 `json:"cargo_vol_behind_2nd_row_cu_ft"`
}

var validBodyStyles = map[string]bool{
	"sedan": true, "hatchback": true, "suv": true, "truck": true,
	"coupe": true, "convertible": true, "wagon": true, "van": true,
	"other": true,
}

var validPowertrains = map[string]bool{
	"gas": true, "diesel": true, "hybrid": true, "electric": true,
	"other": true,
}

func CreateCarHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateCarInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Year < 1980 || input.Year > 2100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "model_year must be between 1980 and 2100"})
			return
		}

		if !validBodyStyles[input.BodyStyle] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body_style"})
			return
		}

		if !validPowertrains[input.Powertrain] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid powertrain"})
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
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{"error": "Car model already exists for this make and model year"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, car)
	}
}

func GetAllCarsHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedQueryParams := map[string]bool{
			"year": true, "make": true, "model": true, "trim": true,
			"body_style": true, "powertrain": true,
		}
		for key := range c.Request.URL.Query() {
			if !allowedQueryParams[key] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported query parameter: " + key})
				return
			}
		}

		var filter FilterRequest
		if err := c.ShouldBindQuery(&filter); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if filter.Year != nil && (*filter.Year < 1980 || *filter.Year > 2100) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "year must be between 1980 and 2100"})
			return
		}

		repositoryFilter := repository.CarFilter{
			Year:       filter.Year,
			Make:       filter.Make,
			Model:      filter.Model,
			Trim:       filter.Trim,
			BodyStyle:  filter.BodyStyle,
			Powertrain: filter.Powertrain,
		}

		cars, err := repository.GetCarsByFilter(pool, repositoryFilter)
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

		var input UpdateCarVolume

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.CargoFullVolumeCuFt == nil && input.CargoBehind2ndRowCuFt == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "At least one cargo volume value is required",
			})
			return
		}

		if (input.CargoFullVolumeCuFt != nil && *input.CargoFullVolumeCuFt < 0) ||
			(input.CargoBehind2ndRowCuFt != nil && *input.CargoBehind2ndRowCuFt < 0) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cargo volume values must be non-negative"})
			return
		}

		car, err := repository.UpdateCarVolume(
			pool,
			id,
			input.CargoBehind2ndRowCuFt,
			input.CargoFullVolumeCuFt,
		)
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

func DeleteCarHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid car ID"})
			return
		}

		err = repository.DeleteCar(pool, id)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
				return
			}

			//TODO handle status 404 not found
			if err.Error() == fmt.Sprintf("car with id %d not found", id) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Car deleted successfully"})
	}
}
