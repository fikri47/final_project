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

func GetSocialMedias(c *gin.Context) {
	var (
		db              = database.GetDB()
		SocialMedias    = []models.SocialMedia{}
		GetSocialMedias = []models.GetSocialMedias{}
	)

	err := db.Preload("User").Find(&SocialMedias).Error

	if err != nil {
		writeError(c, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, socialMedia := range SocialMedias {
		GetSocialMedias = append(GetSocialMedias, models.GetSocialMedias{
			Id:             socialMedia.Id,
			Name:           socialMedia.Name,
			SocialMediaUrl: socialMedia.SocialMediaUrl,
			UserId:         socialMedia.UserId,
			CreatedAt:      socialMedia.CreatedAt,
			UpdateAt:       socialMedia.UpdatedAt,
			User: models.UserSocialMedia{
				Id:       socialMedia.User.Id,
				Username: socialMedia.User.Username,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"social_medias": GetSocialMedias,
	})
}

func CreateSocialMedia(c *gin.Context) {
	var (
		db          = database.GetDB()
		userData    = c.MustGet("userData").(jwt.MapClaims)
		contentType = helpers.GetContentType(c)
		SocialMedia = models.SocialMedia{}
		userId      = uint(userData["id"].(float64))
		err         error
	)

	if contentType == appJSON {
		err = c.ShouldBindJSON(&SocialMedia)
	} else {
		err = c.ShouldBind(&SocialMedia)
	}

	if err != nil {
		writeError(c, "Internal server error", http.StatusInternalServerError)
		return
	}

	if SocialMedia.Name == "" {
		writeError(c, "Your social media name is required", http.StatusBadRequest)
		return
	}

	if SocialMedia.SocialMediaUrl == "" {
		writeError(c, "Your social media url is required", http.StatusBadRequest)
		return
	}

	SocialMedia.UserId = userId

	err = db.Debug().Create(&SocialMedia).Error

	if err != nil {
		writeError(c, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":               SocialMedia.Id,
		"name":             SocialMedia.Name,
		"social_media_url": SocialMedia.SocialMediaUrl,
		"user_id":          SocialMedia.UserId,
		"created_at":       SocialMedia.CreatedAt,
	})

}

func UpdateSocialMedia(c *gin.Context) {
	var (
		db                 = database.GetDB()
		userData           = c.MustGet("userData").(jwt.MapClaims)
		contentType        = helpers.GetContentType(c)
		SocialMedia        = models.SocialMedia{}
		socialMediaId, err = strconv.Atoi(c.Param("socialMediaId"))
		userId             = uint(userData["id"].(float64))
	)

	if err != nil {
		writeError(c, "invalid Paramater", http.StatusBadRequest)
		return
	}

	if contentType == appJSON {
		err = c.ShouldBindJSON(&SocialMedia)
	} else {
		err = c.ShouldBind(&SocialMedia)
	}

	if err != nil {
		writeError(c, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if SocialMedia.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Your social media name is required",
		})
		return
	}

	if SocialMedia.SocialMediaUrl == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Your social media url is required",
		})
		return
	}

	err = db.Select("user_id").First(&SocialMedia, socialMediaId).Error

	if err != nil {
		writeError(c, "Data Doesn't Exist", http.StatusNotFound)
		return
	}

	if SocialMedia.UserId != userId {
		writeError(c, "You are not allowed to access this data", http.StatusUnauthorized)
		return
	}

	SocialMedia.UserId = userId
	SocialMedia.Id = uint(socialMediaId)

	err = db.Model(&SocialMedia).Where("id=?", socialMediaId).Updates(SocialMedia).Error

	if err != nil {
		writeError(c, "Data Doest's Exist", http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":               SocialMedia.Id,
		"name":             SocialMedia.Name,
		"social_media_url": SocialMedia.SocialMediaUrl,
		"user_id":          SocialMedia.UserId,
		"updated_at":       SocialMedia.UpdatedAt,
	})
}

func DeleteSocialMedia(c *gin.Context) {
	var (
		db                 = database.GetDB()
		SocialMedia        = models.SocialMedia{}
		userData           = c.MustGet("userData").(jwt.MapClaims)
		userId             = uint(userData["id"].(float64))
		socialMediaId, err = strconv.Atoi(c.Param("socialMediaId"))
	)

	if err != nil {
		writeError(c, "Invalid parameter", http.StatusBadRequest)
		return
	}

	err = db.Select("user_id").First(&SocialMedia, socialMediaId).Error

	if err != nil {
		writeError(c, "Data Doesn't Exist", http.StatusNotFound)
		return
	}

	if SocialMedia.UserId != userId {
		writeError(c, "You are not allowed to access this data", http.StatusUnauthorized)
		return
	}

	err = db.Delete(models.SocialMedia{}, "id", socialMediaId).Error

	if err != nil {
		writeError(c, "Error Deleting Item", http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your social media has been successfully deleted",
	})
}
