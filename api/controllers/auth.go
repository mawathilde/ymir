package controllers

import (
	"net/http"
	"time"
	"ymir/api/db"
	"ymir/api/models"
	"ymir/api/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	// Get email & pass off req body
	var loginRequest models.LoginRequest

	if c.Bind(&loginRequest) != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "Failed to read body",
		})

		return
	}

	// Look up for requested user
	var user models.User

	db.DB.First(&user, "username = ?", loginRequest.Username)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "Invalid username or password",
		})
		return
	}

	// Compare sent in password with saved users password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "Invalid email or password",
		})
		return
	}

	// Generate a JWT token
	tokenString, err := utils.GenerateToken(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// Respond
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24, "", "", false, true)

	tokenResponse := models.TokenResponse{Token: tokenString, ExpiresAt: time.Now().Add(time.Hour * 24).Unix()}

	c.JSON(http.StatusOK, tokenResponse)
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": "You are authenticated!",
		"user":    user,
	})
}
