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

	c.SetCookie("user_id", userInfo.UserId, 0, "/", "localhost:3000", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func GetUser(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	userId := c.Query("user_id")

	user, err := models.GetUserInfo("user_id", userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "User not found",
		})
		return
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
