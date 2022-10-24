package controller

import (
	"final_project/database"
	"final_project/helpers"
	"final_project/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GetPhoto(c *gin.Context) {
	var (
		db        = database.GetDB()
		Photo     = []models.Photo{}
		GetPhotos = []models.GetPhotos{}
		err       error
	)

	err = db.Preload("User").Find(&Photo).Error

	if err != nil {
		writeError(c, "Internal Server error", http.StatusInternalServerError)
		return
	}

	for _, photo := range Photo {
		GetPhotos = append(GetPhotos, models.GetPhotos{
			Id:       photo.Id,
			Title:    photo.Title,
			Caption:  photo.Caption,
			PhotoUrl: photo.PhotoUrl,
			UserId:   photo.UserId,
			CreateAt: photo.CreatedAt,
			UpdateAt: photo.UpdatedAt,
			User: models.UserPhoto{
				Email:    photo.User.Email,
				Username: photo.User.Username,
			},
		})
	}

	c.JSON(http.StatusOK, GetPhotos)
}

func CreatePhoto(c *gin.Context) {
	var (
		db          = database.GetDB()
		userData    = c.MustGet("userData").(jwt.MapClaims)
		userId      = uint(userData["id"].(float64))
		contentType = helpers.GetContentType(c)
		Photo       = models.Photo{}
		err         error
	)

	if contentType == appJSON {
		err = c.ShouldBindJSON(&Photo)
	} else {
		err = c.ShouldBind(&Photo)
	}

	if err != nil {
		writeError(c, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if Photo.Title == "" {
		writeError(c, "Your photo title is required", http.StatusBadRequest)
		return
	}

	if Photo.PhotoUrl == "" {
		writeError(c, "your photo_url is required", http.StatusBadRequest)
		return
	}

	Photo.UserId = uint(userId)

	err = db.Debug().Create(&Photo).Error

	if err != nil {
		writeError(c, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         Photo.Id,
		"title":      Photo.Title,
		"caption":    Photo.Caption,
		"photo_url":  Photo.PhotoUrl,
		"user_id":    Photo.UserId,
		"created_at": Photo.CreatedAt,
	})
}

func UpdatePhoto(c *gin.Context) {
	var (
		db           = database.GetDB()
		contentType  = helpers.GetContentType(c)
		userData     = c.MustGet("userData").(jwt.MapClaims)
		userId       = uint(userData["id"].(float64))
		photoId, err = strconv.Atoi(c.Param("photoId"))
		Photo        = models.Photo{}
	)

	if err != nil {
		writeError(c, "Invalid Parameter", http.StatusBadRequest)
		return
	}

	if contentType == appJSON {
		err = c.ShouldBindJSON(&Photo)
	} else {
		err = c.ShouldBind(&Photo)
	}

	if err != nil {
		writeError(c, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if Photo.Title == "" {
		writeError(c, "Your photo title is required", http.StatusBadRequest)
		return
	}

	if Photo.PhotoUrl == "" {
		writeError(c, "Your photo title is required", http.StatusBadRequest)
		return
	}

	err = db.Select("user_id").First(&Photo, photoId).Error

	if err != nil {
		writeError(c, "Data Doesn't Exist", http.StatusNotFound)
		return
	}

	if Photo.UserId != userId {
		writeError(c, "You are not allowed to access this data", http.StatusUnauthorized)
		return
	}

	Photo.Id = uint(photoId)

	err = db.Debug().Model(&Photo).Where("id=?", photoId).Updates(&Photo).Error

	if err != nil {
		writeError(c, "Data Doesn't Exist", http.StatusNotFound)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         Photo.Id,
		"title":      Photo.Title,
		"caption":    Photo.Caption,
		"photo_url":  Photo.PhotoUrl,
		"user_id":    Photo.UserId,
		"updated_at": Photo.UpdatedAt,
	})

}

func DeletePhoto(c *gin.Context) {
	var (
		db           = database.GetDB()
		userData     = c.MustGet("userData").(jwt.MapClaims)
		userId       = uint(userData["id"].(float64))
		photoId, err = strconv.Atoi(c.Param("photoId"))
		Photo        = models.Photo{}
	)

	if err != nil {
		writeError(c, "Invalid Paramater", http.StatusBadRequest)
		return
	}

	err = db.Select("user_id").First(&Photo, photoId).Error

	if err != nil {
		writeError(c, "Data Doesn't Exist", http.StatusNotFound)
		return
	}

	if Photo.UserId != userId {
		writeError(c, "You are not allowed to access this data", http.StatusUnauthorized)
		return
	}

	err = db.Delete(models.Comment{}, "photo_id", photoId).Error

	if err != nil {
		writeError(c, "Error Deleting Item", http.StatusInternalServerError)
		return
	}

	err = db.Delete(&Photo, "id", photoId).Error

	if err != nil {
		writeError(c, "Error Deleting Item", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your Photo has been succesfully deleted",
	})

}
