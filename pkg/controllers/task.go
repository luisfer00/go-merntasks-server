package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luisfer00/go-merntasks-server/pkg/models"
	"github.com/luisfer00/go-merntasks-server/pkg/services"
	"github.com/luisfer00/go-merntasks-server/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type createTaskRequestBody struct {
	Nombre string `json:"nombre" binding:"required"`
	Proyecto string `json:"proyecto" binding:"required"`
}

type updateTaskRequestBody struct {
	Nombre *string `json:"nombre"`
	Proyecto string `json:"proyecto" binding:"required"`
	Estado *bool `json:"estado"`
}

func GetTasksController(c *gin.Context) {
	userCtx, _ := c.Get("user")
	userJWT := userCtx.(utils.UserJWT)
	userID, err := primitive.ObjectIDFromHex(userJWT.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}
	projectIDString, ok := c.Params.Get("projectID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": "error with params",
		})
		return
	}

	projectID, err := primitive.ObjectIDFromHex(projectIDString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}

	project, err := services.GetProject(projectID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": true,
				"msj": "project was not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"msj": err.Error(),
			})
		}
		return
	}
	if *project.Creador != userID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": true,
			"msj": "you are not the owner of the project",
		})
		return
	}

	tasks, err := services.GetTasks(project.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func CreateTaskController(c *gin.Context) {
	body := createTaskRequestBody{}
	userCtx, _ := c.Get("user")
	userJWT := userCtx.(utils.UserJWT)
	userID, err := primitive.ObjectIDFromHex(userJWT.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}
	projectID, err := primitive.ObjectIDFromHex(body.Proyecto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}
	
	project, err := services.GetProject(projectID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": true,
				"msj": "project was not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"msj": err.Error(),
			})
		}
		return
	}
	if *project.Creador != userID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": true,
			"msj": "you are not the owner of the project",
		})
		return
	}

	taskStatus := false

	task, err := services.InsertTask(models.Task{
		Nombre: body.Nombre,
		Estado: &taskStatus,
		Proyecto: &projectID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func UpdateTaskController(c *gin.Context) {
	body := updateTaskRequestBody{}
	userCtx, _ := c.Get("user")
	userJWT := userCtx.(utils.UserJWT)
	userID, err := primitive.ObjectIDFromHex(userJWT.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}

	projectID, err := primitive.ObjectIDFromHex(body.Proyecto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}
	
	project, err := services.GetProject(projectID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": true,
				"msj": "project was not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"msj": err.Error(),
			})
		}
		return
	}
	if *project.Creador != userID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": true,
			"msj": "you are not the owner of the project",
		})
		return
	}
	
	idParam, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": "error with params",
		})
		return
	}
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}

	task, err := services.GetTask(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": true,
				"msj": "task was not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"msj": err.Error(),
			})
		}
		return
	}
	
	if body.Nombre != nil && *body.Nombre != "" {
		task.Nombre = *body.Nombre
	}
	if body.Estado != nil {
		task.Estado = body.Estado
	}

	task, err = services.UpdateTask(*task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func DeleteTaskController(c *gin.Context) {
	userCtx, _ := c.Get("user")
	userJWT := userCtx.(utils.UserJWT)
	userID, err := primitive.ObjectIDFromHex(userJWT.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}
	projectIDString, ok := c.Params.Get("projectID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": "error with params",
		})
		return
	}
	idString, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": "error with params",
		})
		return
	}

	projectID, err := primitive.ObjectIDFromHex(projectIDString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}
	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}

	project, err := services.GetProject(projectID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": true,
				"msj": "project was not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"msj": err.Error(),
			})
		}
		return
	}
	if *project.Creador != userID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": true,
			"msj": "you are not the owner of the project",
		})
		return
	}

	services.DeleteTask(id)
	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"msj": "task has been deleted successfully",
	})
}

func DeleteTasksController(c *gin.Context) {
	userCtx, _ := c.Get("user")
	userJWT := userCtx.(utils.UserJWT)
	userID, err := primitive.ObjectIDFromHex(userJWT.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}
	projectIDString, ok := c.Params.Get("projectID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": "error with params",
		})
		return
	}

	projectID, err := primitive.ObjectIDFromHex(projectIDString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}

	project, err := services.GetProject(projectID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": true,
				"msj": "project was not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"msj": err.Error(),
			})
		}
		return
	}
	if *project.Creador != userID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": true,
			"msj": "you are not the owner of the project",
		})
		return
	}

	services.DeleteTasks(projectID)
	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"msj": "tasks have been deleted successfully",
	})
}