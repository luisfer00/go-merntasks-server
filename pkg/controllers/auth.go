package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/luisfer00/go-merntasks-server/pkg/services"
	"github.com/luisfer00/go-merntasks-server/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type loginUserRequestBody struct {
	Email string `json:"email" binding:"email"`
	Password string `json:"password" binding:"gte=6"`
}

func LoginUserController(c *gin.Context) {
	SECRET := os.Getenv("SECRET")
	body := loginUserRequestBody{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msj": "Error parsing json",
		})
		return
	}

	existingUser, err := services.GetUser(body.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": true,
				"msj": "user was not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"msj": err.Error(),
			})
		}
		return
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(body.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": true,
				"msj": "wrong password",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"msj": err.Error(),
			})
			return
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usuario": jwt.MapClaims{
			"id": existingUser.ID,
		},
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	jwt, err := token.SignedString([]byte(SECRET))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"token": jwt,
	})
}

func GetUserDataController(c *gin.Context) {
	userCtx, _ := c.Get("user")
	userJWT, _  := userCtx.(utils.UserJWT)
	userID, err := primitive.ObjectIDFromHex(userJWT.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}

	user, err := services.GetUserByID(userID)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": "user doesnt exist",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}

	user.Password = ""

	c.JSON(200, *user)
}