package utils

import "github.com/golang-jwt/jwt/v4"

type AuthClaims struct {
	User UserJWT `json:"usuario"`
	jwt.RegisteredClaims
}

type UserJWT struct {
	Id string `json:"id"`
}