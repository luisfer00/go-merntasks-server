package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/luisfer00/go-merntasks-server/pkg/controllers"
	"github.com/luisfer00/go-merntasks-server/pkg/middlewares"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("")
	authGroup.Use(middlewares.Auth())

	rg.POST("", controllers.LoginUserController)
	authGroup.GET("", controllers.GetUserDataController)
}