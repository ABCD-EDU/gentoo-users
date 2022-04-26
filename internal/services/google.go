package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/abcd-edu/gentoo-users/internal/models"
	helpers "github.com/abcd-edu/gentoo-users/internal/services/helpers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfGl = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:8001/callback-gl",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateStringGl = ""
)

func InitializeOAuthGoogle() {
	oauthConfGl.ClientID = viper.GetString("google.clientID")
	oauthConfGl.ClientSecret = viper.GetString("google.clientSecret")
	oauthStateStringGl = viper.GetString("oauthStateStringGl")
}

func HandleGoogleLogin(c *gin.Context) {
	HandleLogin(c, oauthConfGl, oauthStateStringGl)
}

/*
	This validates if the user is properly loggged in with Google.
	Below are the steps in verification and where to proceed after validating
	state, tokens, etc.

	1. Check if user is authenticated through "verified_email" callback response
	- If user is not authenticated
		- Redirect back to "/"
	- If user is authenticated
		- Proceed to check if user is registered

	2. Check if user is registered on the database via email
	- If user is Registered:
		- Redirect back to "/"
	- If user is NOT Registered
		- Redirect content to "/register-gl"
*/
func CallBackFromGoogle(c *gin.Context) {
	if c.Request.FormValue("state") != oauthStateStringGl {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	token, err := oauthConfGl.Exchange(oauth2.NoContext, c.Request.FormValue("code"))
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	resBytes := []byte(string(string(content)))
	var jsonRes map[string]interface{}
	_ = json.Unmarshal(resBytes, &jsonRes)
	verified_email := jsonRes["verified_email"].(bool)

	if verified_email {
		email := jsonRes["verified_email"].(string)
		helpers.SetAuthenticationKey(c)
		helpers.SetUserInfo(email, c)

		_, err := models.GetUserInfo(email)
		if err == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/register-gl")
		}
	}
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
