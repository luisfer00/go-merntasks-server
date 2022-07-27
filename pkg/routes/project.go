package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/luisfer00/go-merntasks-server/pkg/controllers"
	"github.com/luisfer00/go-merntasks-server/pkg/middlewares"
)

func RegisterProjectRoutes(rg *gin.RouterGroup) {
	rg.Use(middlewares.Auth())
	rg.GET("", controllers.GetProjectsController)
	rg.POST("", controllers.CreateProjectController)
	rg.PUT("/:id", controllers.UpdateProjectController)
	rg.DELETE("/:id", controllers.DeleteProjectController)
}