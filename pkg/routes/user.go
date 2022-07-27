package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/luisfer00/go-merntasks-server/pkg/controllers"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {
	rg.POST("", controllers.RegisterUserController)
}