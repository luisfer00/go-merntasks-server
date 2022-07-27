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

type createProjectRequestBody struct {
	Nombre string `json:"nombre" binding:"required"`
}
type updateProjectRequestBody struct {
	Nombre string `json:"nombre" binding:"required"`
}

func GetProjectsController(c *gin.Context) {
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

	projects, err := services.GetProjects(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func CreateProjectController(c *gin.Context) {
	body := createProjectRequestBody{}
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

	project, err := services.InsertProject(models.Project{
		Nombre: body.Nombre,
		Creador: &userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
	}

	c.JSON(http.StatusOK, project)
}

func UpdateProjectController(c *gin.Context) {
	body := updateProjectRequestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msj": "Error parsing json",
		})
		return
	}
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

	project, err := services.GetProject(id)
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
	if body.Nombre == project.Nombre {
		c.JSON(http.StatusOK, project)
		return
	}
	if body.Nombre != "" {
		project.Nombre = body.Nombre
	}

	project, err = services.UpdateProject(*project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, project)
}

func DeleteProjectController(c *gin.Context) {
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
	project, err := services.GetProject(id)
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
	
	services.DeleteProject(id)
	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"msj": "project has been deleted successfully",
	})
}