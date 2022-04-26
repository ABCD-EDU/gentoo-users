package main

import (
	"github.com/abcd-edu/gentoo-users/internal/configs"
	"github.com/abcd-edu/gentoo-users/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	configs.InitializeViper()

	services.InitializeOAuthGoogle()

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.GET("/", services.HandleMain)
		v1.POST("/login-gl", services.HandleGoogleLogin)
		v1.POST("/callback-gl", services.CallBackFromGoogle)
	}

	port := viper.GetString("port")
	router.Run(":" + port)
}
