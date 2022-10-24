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

func GetComment(c *gin.Context) {
	var (
		db          = database.GetDB()
		Comment     = []models.Comment{}
		GetComments = []models.GetComment{}
		err         error
	)
	err = db.Preload("User").Preload("Photo").Find(&Comment).Error

	if err != nil {
		writeError(c, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, comment := range Comment {
		GetComments = append(GetComments, models.GetComment{
			Id:        comment.Id,
			Message:   comment.Message,
			PhotoId:   comment.PhotoId,
			UserId:    comment.UserId,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			User: models.UserComment{
				Id:       comment.User.Id,
				Email:    comment.User.Email,
				Username: comment.User.Username,
			},
			Photo: models.PhotoComment{
				Id:       comment.Photo.Id,
				Title:    comment.Photo.Title,
				Caption:  comment.Photo.Caption,
				PhotoUrl: comment.Photo.PhotoUrl,
				UserId:   comment.Photo.UserId,
			},
		})
	}

	c.JSON(http.StatusOK, GetComments)
}

func CreateComment(c *gin.Context) {
	var (
		db          = database.GetDB()
		userData    = c.MustGet("userData").(jwt.MapClaims)
		userId      = uint(userData["id"].(float64))
		contentType = helpers.GetContentType(c)
		Comment     = models.Comment{}
		NewComment  = models.Comment{}
		err         error
	)

	if contentType == appJSON {
		err = c.ShouldBindJSON(&Comment)
	} else {
		err = c.ShouldBind(&Comment)
	}

	if err != nil {
		writeError(c, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if Comment.Message == "" {
		writeError(c, "Your message is required", http.StatusBadRequest)
		return
	}

	if Comment.PhotoId == 0 {
		writeError(c, "Your Photo ID is required", http.StatusBadRequest)
		return
	}

	err = db.Model(models.Photo{}).Select("id").First(&NewComment, Comment.PhotoId).Error

	if err != nil {
		writeError(c, "Data Doesn't Exist", http.StatusNotFound)
		return
	}

	Comment.UserId = userId

	err = db.Debug().Create(&Comment).Error

	if err != nil {
		writeError(c, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         Comment.Id,
		"message":    Comment.Message,
		"photo_id":   Comment.PhotoId,
		"user_id":    Comment.UserId,
		"created_at": Comment.CreatedAt,
	})
}

func UpdateComment(c *gin.Context) {
	var (
		db             = database.GetDB()
		userData       = c.MustGet("userData").(jwt.MapClaims)
		userId         = uint(userData["id"].(float64))
		contentType    = helpers.GetContentType(c)
		Comment        = models.Comment{}
		commentId, err = strconv.Atoi(c.Param("commentId"))
	)

	if err != nil {
		writeError(c, "Invalid Parameter", http.StatusBadRequest)
		return
	}

	if contentType == appJSON {
		err = c.ShouldBindJSON(&Comment)
	} else {
		err = c.ShouldBind(&Comment)
	}

	if err != nil {
		writeError(c, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if Comment.Message == "" {
		writeError(c, "Your Message is required", http.StatusBadRequest)
		return
	}

	err = db.Select("user_id", "photo_id").First(&Comment, commentId).Error

	if err != nil {
		writeError(c, "Data Doesn't Exist", http.StatusNotFound)
		return
	}

	if Comment.UserId != userId {
		writeError(c, "You are not allowed to access this data", http.StatusUnauthorized)
		return
	}

	Comment.Id = uint(commentId)

	err = db.Debug().Model(&Comment).Where("id=?", commentId).Updates(&Comment).Error

	if err != nil {
		writeError(c, "Data Doesn't Exist", http.StatusNotFound)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         Comment.Id,
		"message":    Comment.Message,
		"photo_id":   Comment.PhotoId,
		"user_id":    Comment.UserId,
		"updated_at": Comment.UpdatedAt,
	})

}

func DeleteComment(c *gin.Context) {
	var (
		db             = database.GetDB()
		userData       = c.MustGet("userData").(jwt.MapClaims)
		userId         = uint(userData["id"].(float64))
		Comment        = models.Comment{}
		commentId, err = strconv.Atoi(c.Param("commentId"))
	)

	if err != nil {
		writeError(c, "Invalid Parameter", http.StatusBadRequest)
		return
	}

	err = db.Select("user_id").First(&Comment, commentId).Error

	if err != nil {
		writeError(c, "Data Doesn't Exist", http.StatusNotFound)
		return
	}

	if Comment.UserId != userId {
		writeError(c, "You are not allowed to access this data", http.StatusUnauthorized)
		return
	}

	err = db.Delete(&Comment, "id", commentId).Error

	if err != nil {
		writeError(c, "Error Deleting item", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your Comment Has been Succesfully deleted",
	})
}
