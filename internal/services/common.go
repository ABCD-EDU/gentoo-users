package services

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	helpers "github.com/abcd-edu/gentoo-users/internal/services/helpers"
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
	helpers.RemoveCookies(c)
}
