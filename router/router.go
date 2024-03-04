package router

import (
	"api-golang/controllers"
	"api-golang/middlewares"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.POST("api/signup", controllers.SignUp)
	r.POST("api/login", controllers.Login)

	userRoutes := r.Group("/api/users").Use(middlewares.JWTMiddleware())
	{
		userRoutes.GET("/profile", controllers.ProfileUser)
		userRoutes.PUT("/:id", controllers.UpdateUser)
		userRoutes.DELETE("/:id", controllers.DeleteUser)
		userRoutes.POST("/logout", controllers.Logout)
	}

	photoRoutes := r.Group("/api/photos").Use(middlewares.JWTMiddleware())
	{
		photoRoutes.GET("/", controllers.GetPhotos)
		photoRoutes.PUT("/:id", controllers.UpdatePhoto)
		photoRoutes.DELETE("/:id", controllers.DeletePhoto)
	}

	return r
}
