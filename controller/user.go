package controller

import (
	"UnitUser/database"
	"UnitUser/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UserSignUp(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Failed to bind user:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if user.Name == "" || user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "all fields are mandatory"})
		return
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash passowrd"})
		return
	}
	user.Password = string(pass)
	if err := database.DB.Create(&user).Error; err != nil {
		log.Println("Failed to create user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user created successfully."})
}

func UserSignIn(c *gin.Context) {
	var input models.User
	var user models.User

	c.ShouldBindJSON(&input)
	if err := database.DB.Where("email=?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong email or password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}

func UserShow(c *gin.Context) {

	var user []models.User

	if err := database.DB.Find(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": user,
	})
}

func EditUser(c *gin.Context) {
	var user models.User

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}
	updates := map[string]interface{}{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	}
	if err := database.DB.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Printf("Failed to update user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully updated user",
	})
}
