package controller

import (
	"final_project/database"
	"final_project/helpers"
	"final_project/models"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	appJSON = "application/json"
)

func writeError(c *gin.Context, message string, code int) {
	c.JSON(code, gin.H{
		"error_messsage": message,
		"status_code":    code,
	})
}

func UserRegister(c *gin.Context) {
	var (
		db          = database.GetDB()
		contentType = helpers.GetContentType(c)
		User        = models.User{}
		NewUser     = models.User{}
		err         error
	)

	if contentType == appJSON {
		err = c.ShouldBindJSON(&User)
	} else {
		err = c.ShouldBind(&User)
	}

	if err != nil {
		writeError(c, "Internal Server Erro", http.StatusInternalServerError)
		return
	}

	if User.Email == "" {
		writeError(c, "Your Email is Required", http.StatusBadRequest)
		return
	}

	_, errEmail := mail.ParseAddress(User.Email)
	if errEmail != nil {
		writeError(c, "Invalid Email Format", http.StatusBadRequest)
		return
	}

	db.Where("email=?", User.Email).First(&NewUser)

	if NewUser.Email == User.Email {
		writeError(c, "Email Already Used", http.StatusBadRequest)
		return
	}

	if User.Username == "" {
		writeError(c, "Your Username is required", http.StatusBadRequest)
		return
	}

	db.Where("username", User.Username).First(&NewUser)

	if NewUser.Username == User.Username {
		writeError(c, "Username Already User", http.StatusBadRequest)
		return
	}

	if User.Password == "" {
		writeError(c, "Your password is required", http.StatusBadRequest)
		return
	}

	if len(User.Password) < 6 {
		writeError(c, "Password has to have a minimum length of 6 characters", http.StatusBadRequest)
		return
	}

	if User.Age == 0 {
		writeError(c, "Your age is required", http.StatusBadRequest)
		return
	}

	if User.Age <= 8 {
		writeError(c, "Sorry, you must be at least 8 years old", http.StatusBadRequest)
		return
	}

	err = db.Debug().Create(&User).Error

	if err != nil {
		writeError(c, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"age":      User.Age,
		"email":    User.Email,
		"id":       User.Id,
		"username": User.Username,
	})
}

func UserLogin(c *gin.Context) {
	var (
		db          = database.GetDB()
		contentType = helpers.GetContentType(c)
		User        = models.User{}
		password    = ""
		err         error
	)

	if contentType == appJSON {
		err = c.ShouldBindJSON(&User)
	} else {
		err = c.ShouldBind(&User)
	}

	if err != nil {
		writeError(c, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	password = User.Password

	err = db.Debug().Where("email=?", User.Email).Take(&User).Error
	comparePass := helpers.ComparePass([]byte(User.Password), []byte(password))

	if err != nil || !comparePass {
		writeError(c, "Invalid Username or password", http.StatusUnauthorized)
		return
	}

	token, err := helpers.GenerateToken(
		User.Id,
		User.Email,
	)

	if err != nil {
		writeError(c, "Unauthorized", http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func UserUpdate(c *gin.Context) {
	var (
		db          = database.GetDB()
		userData    = c.MustGet("userData").(jwt.MapClaims)
		contentType = helpers.GetContentType(c)
		User        = models.User{}
		NewUser     = models.User{}
		userId      = userData["id"].(float64)
		err         error
	)

	paramUserId, err := strconv.Atoi(c.Param("userId"))

	if err != nil {
		writeError(c, "Invalid Parameter", http.StatusBadRequest)
		return
	}

	if contentType == appJSON {
		err = c.ShouldBindJSON(&User)
	} else {
		err = c.ShouldBind(&User)
	}

	if err != nil {
		writeError(c, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if User.Email == "" {
		writeError(c, "Your Email is Required", http.StatusBadRequest)
		return
	}

	_, errEmail := mail.ParseAddress(User.Email)
	if errEmail != nil {
		writeError(c, "Invalid Email Format", http.StatusBadRequest)
		return
	}

	db.Where("email=?", User.Email).First(&NewUser)

	if NewUser.Email == User.Email {
		writeError(c, "Email Already Used", http.StatusBadRequest)
		return
	}

	if User.Username == "" {
		writeError(c, "Your Username is Required", http.StatusBadRequest)
		return
	}

	db.Where("username", User.Username).First(&NewUser)

	if NewUser.Username == User.Username {
		writeError(c, "Username Already Used", http.StatusBadRequest)
		return
	}

	err = db.Select("id", "age").First(&User, paramUserId).Error

	if err != nil {
		writeError(c, "Data Not Found", http.StatusNotFound)
		return
	}

	if User.Id != uint(userId) {
		writeError(c, "You Are not Allowed to access this data", http.StatusUnauthorized)
		return
	}

	db.Where("email=?", User.Email).First(&NewUser)

	if User.Email == NewUser.Email {
		writeError(c, "Email Already Used", http.StatusBadRequest)
		return
	}

	db.Where("username=?", User.Username).First(&NewUser)

	if User.Username == NewUser.Username {
		writeError(c, "Username Already Used", http.StatusBadRequest)
		return
	}

	err = db.Model(&User).Where("id=?", paramUserId).Updates(&User).Error

	if err != nil {
		writeError(c, "Data Not Found", http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         User.Id,
		"email":      User.Email,
		"username":   User.Username,
		"age":        User.Age,
		"updated_at": User.UpdatedAt,
	})

}

func UserDelete(c *gin.Context) {
	var (
		db               = database.GetDB()
		userData         = c.MustGet("userData").(jwt.MapClaims)
		userId           = uint(userData["id"].(float64))
		paramUserId, err = strconv.Atoi(c.Param("userId"))
		User             = models.User{}
	)

	if err != nil {
		writeError(c, "Invalid Parameter", http.StatusBadRequest)
		return
	}

	err = db.Select("id").First(&User, paramUserId).Error

	if err != nil {
		writeError(c, "Data Not Found", http.StatusNotFound)
		return
	}

	if userId != uint(paramUserId) {
		writeError(c, "You are not allowed to access this data", http.StatusUnauthorized)
		return
	}

	err = db.Delete(models.SocialMedia{}, "user_id", userId).Error

	if err != nil {
		writeError(c, "Error Deleting Item", http.StatusInternalServerError)
		return
	}

	err = db.Delete(models.Comment{}, "user_id", userId).Error

	if err != nil {
		writeError(c, "Error Deleting Item", http.StatusInternalServerError)
		return
	}

	err = db.Delete(models.Photo{}, "user_id", userId).Error

	if err != nil {
		writeError(c, "Error Deleting Item", http.StatusInternalServerError)
		return
	}

	err = db.Delete(User, "id", paramUserId).Error

	if err != nil {
		writeError(c, "Error Deleting Item", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your Account has been succesfully deleted",
	})
}
