package controllers

import (
	"api-golang/database"
	"api-golang/models"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetPhotos(c *gin.Context) {
	var photos []models.Photo
	database.DB.Find(&photos)

	c.JSON(http.StatusOK, gin.H{"data": photos})
}

func UploadPhoto(c *gin.Context) {
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large. Please upload a file less than 10MB"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not found"})
		return
	}
	defer file.Close()

	var photo models.Photo
	photo.Title = c.Request.FormValue("title")
	photo.Caption = c.Request.FormValue("caption")

	err = SaveImageToDatabase(&photo, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save image to database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Image uploaded successfully",
	})

}

func UpdatePhoto(c *gin.Context) {
	var photo models.Photo
	if err := database.DB.First(&photo, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	if err := c.BindJSON(&photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	database.DB.Save(&photo)
	c.JSON(http.StatusOK, gin.H{"data": photo})
}

func DeletePhoto(c *gin.Context) {
	var photo models.Photo
	if err := database.DB.First(&photo, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	database.DB.Delete(&photo)
	c.JSON(http.StatusOK, gin.H{"data": true})
}

func SaveImageToDatabase(photo *models.Photo, file multipart.File) error {
	// read the file
	imagePath := "image/" + photo.Title + ".jpg"
	outFile, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		return err
	}

	photo.PhotoUrl = imagePath
	return nil
}
