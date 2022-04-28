package helpers

import (
	"github.com/gin-gonic/gin"
)

var (
	day = 1
)

func SetAuthenticationKey(c *gin.Context, authToken string) {
	_, err := c.Cookie("auth_key")

	if err != nil {
		c.SetCookie(
			"auth_key",
			authToken,
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
