package middleware

import (
	"net/http"
	"time"
	"ymir/api/db"
	"ymir/api/models"
	"ymir/api/utils"

	"github.com/gin-gonic/gin"
)

func RequireAuth(c *gin.Context) {
	// Get the cookie off the request
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		// get bearer token from header
		tokenString = c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
			return
		} else {
			// Remove the bearer prefix
			tokenString = tokenString[7:]
		}
	}

	// Decode/validate it
	claims, err := utils.ParseToken(tokenString)

	if err == nil && claims != nil {
		// Check the expiry date
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			return
		}

		// Find the user with token Subject
		var user models.User
		db.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Valid token but user not found, wait what?"})
			return
		}

		// Attach the request
		c.Set("user", user)
		c.Set("claims", claims)

		// Continue
		c.Next()
	} else {
		// Abort with a message
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
	}

}
