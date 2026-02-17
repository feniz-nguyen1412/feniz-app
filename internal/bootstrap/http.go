package bootstrap

import (
	"gomono_template/internal/user/interfaces/http"
	"github.com/gin-gonic/gin"
)

func NewHTTPServer(container *Container) *gin.Engine {
	r := gin.Default()

	http.SetupUserRoutes(r, container.UserHandler)

	return r
}
