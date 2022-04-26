package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	day = 24 * 60 * 60 * 1000
)

func SetAuthenticationKey(c *gin.Context) {
	_, err := c.Cookie("auth_key")

	if err != nil {
		c.SetCookie(
			"auth_key",
			viper.GetString("auth_key"),
			day,
			"/",
			"localhost",
			false,
			true,
		)
	}
}

func SetUserInfo(email string, c *gin.Context) {
	_, err := c.Cookie("email")

	if err != nil {
		c.SetCookie(
			"email",
			email,
			day,
			"/",
			"localhost",
			false,
			true,
		)
	}
}

func RemoveCookies(c *gin.Context) {
	_, err := c.Cookie("email")

	if err != nil {
		c.SetCookie(
			"email",
			"",
			-1,
			"/",
			"localhost",
			false,
			true,
		)
	}

	_, err = c.Cookie("auth_key")

	if err != nil {
		c.SetCookie(
			"auth_key",
			"",
			-1,
			"/",
			"localhost",
			false,
			true,
		)
	}
}
