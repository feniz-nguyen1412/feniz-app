package http

import (
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, handler *UserHandler) {
	users := r.Group("/users")
	{
		users.POST("", handler.CreateUser)
		users.GET("", handler.GetAllUsers)
		users.GET("/:id", handler.GetUser)
		users.PUT("/:id", handler.UpdateUser)
		users.DELETE("/:id", handler.DeleteUser)
	}
}