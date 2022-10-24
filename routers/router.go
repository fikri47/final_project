package routers

import (
	"final_project/controller"
	"final_project/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	r := gin.Default()

	userRouter := r.Group("/users")
	{
		userRouter.POST("/register", controller.UserRegister)
		userRouter.POST("/login", controller.UserLogin)
		userRouter.Use(middlewares.Authentication())
		userRouter.PUT("/:userId", controller.UserUpdate)
		userRouter.DELETE("/:userId", controller.UserDelete)
	}
	socialMediasRouter := r.Group("/socialmedias")
	{
		socialMediasRouter.Use(middlewares.Authentication())
		socialMediasRouter.PUT("/:socialMediaId", controller.UpdateSocialMedia)
		socialMediasRouter.DELETE("/:socialMediaId", controller.DeleteSocialMedia)
		socialMediasRouter.POST("/", controller.CreateSocialMedia)
		socialMediasRouter.GET("/", controller.GetSocialMedias)
	}

	photoRouter := r.Group("/photos")
	{
		photoRouter.Use(middlewares.Authentication())
		photoRouter.POST("/", controller.CreatePhoto)
		photoRouter.GET("/", controller.GetPhoto)
		photoRouter.PUT("/:photoId", controller.UpdatePhoto)
		photoRouter.DELETE("/:photoId", controller.DeletePhoto)
	}

	commentsRouter := r.Group("/comments")
	{
		commentsRouter.Use(middlewares.Authentication())
		commentsRouter.POST("/", controller.CreateComment)
		commentsRouter.GET("/", controller.GetComment)
		commentsRouter.PUT("/:commentId", controller.UpdateComment)
		commentsRouter.DELETE("/:commentId", controller.DeleteComment)
	}

	return r
}
