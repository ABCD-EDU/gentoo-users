package services

import (
	"fmt"
	"net/http"

	"github.com/abcd-edu/gentoo-users/internal/models"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	email := c.Query("email")
	username := c.Query("username")
	google_photo := c.Query("google_photo")
	description := c.Query("description")

	var test models.UserRegistration
	if err := c.ShouldBindJSON(&test); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"RHUIAWHDIUWAHEDUIWAHDUIAWHI": err.Error()})
		return
	}

	fmt.Printf("email: %s\n", test.Email)
	fmt.Printf("usernaem: %s\n", test.Username)
	fmt.Printf("photo: %s\n", test.GooglePhoto)
	fmt.Printf("desc: %s\n", test.Description)

	_ = models.UserRegistration{Email: email, Username: username, GooglePhoto: google_photo, Description: description}

	err := models.RegisterUser(test)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "User is already registered",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func GetUser(c *gin.Context) {
	email, err := c.Cookie("email")
	c.Header("Content-Type", "application/json")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "User not found",
		})
	}

	user, err := models.GetUserInfo(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "User not found",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User found",
		"user":    user,
	})
}

func UpdateUser(c *gin.Context) {

}

func DeleteUser(c *gin.Context) {

}
