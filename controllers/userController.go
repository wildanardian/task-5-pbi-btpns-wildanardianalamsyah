package controllers

import (
	"api-golang/app"
	"api-golang/database"
	"api-golang/helpers"
	"api-golang/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	//get the email and password off request body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	//create a new user
	user := models.User{Email: body.Email, Password: string(hash)}
	result := database.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})

}

func Login(c *gin.Context) {
	// Get the email and password from the request body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, app.CustomError{Code: http.StatusBadRequest, Message: "Invalid request"})
		return
	}

	// Find the user
	var user models.User
	if err := database.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, app.CustomError{Code: http.StatusNotFound, Message: "User not found"})
		return
	}

	// Compare the passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusBadRequest, app.CustomError{Code: http.StatusBadRequest, Message: "Invalid password"})
		return
	}

	// Generate the JWT token
	tokenString, err := helpers.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, app.CustomError{Code: http.StatusInternalServerError, Message: "Failed to generate token"})
		return
	}

	// Set the cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	})

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"status code": http.StatusOK,
		"token":       tokenString,
	})
	c.Set("user_id", user.ID)
	c.Set("user", user)
}

func ProfileUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"user_id": c.GetInt("user_id"),
	})
}

func UpdateUser(c *gin.Context) {
	var user models.User
	if err := database.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	database.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": user,
	})
}

func DeleteUser(c *gin.Context) {
	var user models.User
	if err := database.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	database.DB.Delete(&user)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": true,
	})
}

func Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "Authorization",
		Value:    "",
		Expires:  time.Now(),
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	})
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "User logged out successfully",
	})
}
