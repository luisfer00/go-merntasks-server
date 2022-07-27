package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/luisfer00/go-merntasks-server/pkg/controllers"
	"github.com/luisfer00/go-merntasks-server/pkg/middlewares"
)

func RegisterTaskRoutes(rg *gin.RouterGroup) {
	rg.Use(middlewares.Auth())
	rg.GET("/:projectID", controllers.GetTasksController)
	rg.POST("", controllers.CreateTaskController)
	rg.PUT("/:id", controllers.UpdateTaskController)
	rg.DELETE("/:projectID", controllers.DeleteTasksController)
	rg.DELETE("/:projectID/:id", controllers.DeleteTaskController)
}