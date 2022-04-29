package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/abcd-edu/gentoo-users/internal/models"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user models.UserRegistration
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Content is not valid": err.Error()})
		return
	}

	fmt.Printf("email: %s\n", user.Email)
	fmt.Printf("usernaem: %s\n", user.Username)
	fmt.Printf("photo: %s\n", user.GooglePhoto)
	fmt.Printf("desc: %s\n", user.Description)

	userInfo, err := models.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "User is already registered",
		})
		return
	}

	userJson, err := json.Marshal(userInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "User is already registered",
		})
		return
	}

	fmt.Println(userJson)
	_, err = http.Post("http://localhost:8002/v1/register", "application/json", bytes.NewBuffer(userJson))
	if err != nil {
		fmt.Println("CANNOT SEND POST REQUEST TO POST SERVICE")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Something went wrong with Post Microservice",
		})
		return
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
