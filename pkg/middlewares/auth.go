package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/luisfer00/go-merntasks-server/pkg/utils"
)

var (
	ErrEmptyAuthHeader = errors.New("authorization header not found")
)

func Auth() gin.HandlerFunc {
	SECRET := os.Getenv("SECRET")

	return func(c *gin.Context) {
		tokenString := c.GetHeader("x-auth-token")
		if tokenString == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"msj": "Authorization header not found",
			})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &utils.AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(SECRET), nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"msj": err.Error(),	
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*utils.AuthClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"msj": "error parsing token",	
			})
			c.Abort()
			return
		}
		
		c.Set("user", claims.User)
		c.Next()
	}
}