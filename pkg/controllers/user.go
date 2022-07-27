package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/luisfer00/go-merntasks-server/pkg/models"
	"github.com/luisfer00/go-merntasks-server/pkg/services"
	"github.com/luisfer00/go-merntasks-server/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	SECRET = os.Getenv("SECRET")
)

type registerUserRequestBody struct {
	Nombre string `json:"nombre" binding:"required"`
	Email string `json:"email" binding:"email"`
	Password string `json:"password" binding:"gte=6"`
}

func RegisterUserController(c *gin.Context) {
	body := registerUserRequestBody{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msj": "Error parsing json",
		})
		return
	}

	passwordData, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msj": err.Error(),
		})
		return
	}
	password := string(passwordData)

	newUser, err := services.InsertUser(models.User{
		Nombre: body.Nombre,
		Email: body.Email,
		Password: password,
	})
	if err != nil {
		if err == services.ErrUserExist {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": true,
				"msj": err.Error(),
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

	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.AuthClaims{
		User: utils.UserJWT{
			Id: newUser.ID.Hex(),
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
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