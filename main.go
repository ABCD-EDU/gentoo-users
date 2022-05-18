package main

import (
	"github.com/abcd-edu/gentoo-users/internal/configs"
	"github.com/abcd-edu/gentoo-users/internal/models"
	"github.com/abcd-edu/gentoo-users/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	configs.InitializeViper()

	services.InitializeOAuthGoogle()

	models.InitializeDB()

	router := gin.Default()
	//store := cookie.NewStore([]byte("secret"))
	//router.Use(sessions.Sessions("gentoo", store))

	router.Use(CORSMiddleware())
	v1 := router.Group("/v1")
	v1.Use(CORSMiddleware())
	{
		v1.GET("/", services.HandleMain)
		v1.GET("/login-gl", services.HandleGoogleLogin)
		v1.GET("/callback-gl", services.CallBackFromGoogle)
		v1.POST("/register", services.Register)
		v1.GET("/logout", services.HandleLogout)
		v1.Any("/auth", services.HandleAuthentication)
		v1.GET("/user", services.GetUser)
		v1.POST("/user/mute", services.MuteUser)
		v1.POST("/user/ban", services.BanUser)
	}

	port := viper.GetString("serverPort")
	router.Run(":" + port)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
