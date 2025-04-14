package main

import (
	"time"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {

}

func corsConfig() cors.Config {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowedMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowedHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposedHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	return config
}

func main() {
	r := gin.Default()
	r.Use(cors.New(corsConfig()))

	//api := r.Group("/")
	//api.Use(middleware.RequireAuth)

	//r.POST("auth/register", controllers.Register)
	//r.POST("auth/login", controllers.Login)
	//r.POST("auth/verify", controllers.Verify)

	//api.GET("auth/validate", middleware.RequireAuth, controllers.Validate)

	r.Run()
}
