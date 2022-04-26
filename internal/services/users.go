package services

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/abcd-edu/gentoo-users/internal/models"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var userRegister models.UserRegistration
	info, _ := json.Marshal(userRegister)

	err := c.ShouldBind(userRegister)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"info":   info,
		})
	}

	err = models.RegisterUser(&userRegister)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"info":   info,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"info":   info,
	})
}

func GetUser(c *gin.Context) {

}

func UpdateUser(c *gin.Context) {

}

func DeleteUser(c *gin.Context) {

}
