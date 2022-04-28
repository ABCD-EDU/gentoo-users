package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfGl = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:8001/v1/callback-gl",
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
	gentooHome := "http://localhost:3000/home"
	gentooRegister := "http://localhost:3000/signup"

	if c.Request.FormValue("state") != oauthStateStringGl {
		log.Printf("Error: problem with state\n")
		c.Redirect(http.StatusTemporaryRedirect, gentooHome)
	}

	token, err := oauthConfGl.Exchange(oauth2.NoContext, c.Request.FormValue("code"))
	if err != nil {
		log.Printf("Error: problem with code\n")
		c.Redirect(http.StatusTemporaryRedirect, gentooHome)
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Printf("Error: problem with access token\n")
		c.Redirect(http.StatusTemporaryRedirect, gentooHome)
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error: problem with response\n")
		c.Redirect(http.StatusTemporaryRedirect, gentooHome)
	}

	resBytes := []byte(string(string(content)))
	var jsonRes map[string]interface{}
	_ = json.Unmarshal(resBytes, &jsonRes)
	verified_email := jsonRes["verified_email"].(bool)
	email := jsonRes["email"].(string)

	if verified_email {
		c.SetCookie("auth", token.AccessToken, 0, "/", "localhost:3000", false, true)
		c.SetCookie("email", email, 0, "/", "localhost:3000", false, true)

		// TODO: (MED PRIO) Instead of redirecting directly to signup page,
		// check if user is registered then let's redirect them back to /home
		c.Redirect(http.StatusTemporaryRedirect, gentooRegister)
	}
	c.Redirect(http.StatusTemporaryRedirect, gentooHome)
}
