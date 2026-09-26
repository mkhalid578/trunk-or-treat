package handlers

import (
	"net/http"

	"car-api/internal/models"
	"car-api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type createUserInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input createUserInput
		err := c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(input.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 6 characters long"})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Failed to hash password": err.Error()})
			return
		}

		user := &models.User{
			Email:    input.Email,
			Password: string(hash),
		}
		createdUser, err := repository.CreateUser(pool, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Failed to create user": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, createdUser)
	}
}
