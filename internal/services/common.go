package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/abcd-edu/gentoo-users/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func HandleLogin(c *gin.Context, oauthConf *oauth2.Config, oauthStateString string) {
	URL, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		fmt.Printf("Error in parsing: %s\n", err)
	}

	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)

	URL.RawQuery = parameters.Encode()
	url := URL.String()
	fmt.Println(url)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleMain(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"api": viper.GetString("oauthStateString"),
	})
}

func HandleLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"state":   true,
		"message": "success",
	})
}

func HandleAuthentication(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	authToken, err := c.Cookie("auth")

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{
			"state":   false,
			"message": "unauthorized",
		})
		return
	}
	fmt.Printf("TOKEN FROM COOKIE: %s\n", authToken)

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + authToken)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{
			"state":      false,
			"message":    "unauthorized",
			"token":      authToken,
			"registered": false,
		})

	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{
			"state":      false,
			"message":    "unauthorized",
			"token":      authToken,
			"registered": false,
		})

	}

	resBytes := []byte(string(string(content)))
	var jsonRes map[string]interface{}
	_ = json.Unmarshal(resBytes, &jsonRes)
	verified_email := jsonRes["verified_email"].(bool)

	if verified_email {
		email, err := c.Cookie("email")
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusNotFound, gin.H{
				"state":      true,
				"message":    "authorized",
				"token":      authToken,
				"registered": false,
			})
			return
		}
		fmt.Printf("EMAIL FROM COOKIE: %s\n", email)

		user, err := models.GetUserInfo("email", email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"state":      true,
				"message":    "authorized",
				"token":      authToken,
				"registered": false,
			})
			return
		}
		fmt.Printf("EMAIL FROM DATABASE: %s\n", user.UserInfo.Email)

		if user.UserInfo.Email == email {
			c.JSON(http.StatusOK, gin.H{
				"state":      true,
				"message":    "authorized",
				"token":      authToken,
				"user_id":    user.UserId,
				"registered": true,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"state":      false,
		"message":    "unauthorized",
		"token":      authToken,
		"registered": false,
	})
}
