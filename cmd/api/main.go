package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/luisfer00/go-merntasks-server/pkg/routes"
)

func init() {
	ENV := os.Getenv("ENV")
	if ENV == "dev" {
		godotenv.Load(".env")
	}
}

func main() {
	PORT := os.Getenv("PORT")
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowCredentials: true,
		AllowHeaders: []string{
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"accept",
			"origin",
			"Cache-Control",
			"X-Requested-With",
			"x-auth-token",
		},
		AllowMethods: []string{
			"POST",
			"HEAD",
			"PATCH",
			"OPTIONS",
			"GET",
			"PUT",
			"DELETE",
		},
	}))

	if PORT == "" {
		log.Fatalln("Port not set")
	}

	userGroup := r.Group("/api/usuarios")
	routes.RegisterUserRoutes(userGroup)
	authGroup := r.Group("/api/auth")
	routes.RegisterAuthRoutes(authGroup)
	projectGroup := r.Group("/api/proyectos")
	routes.RegisterProjectRoutes(projectGroup)
	taskGroup := r.Group("/api/tareas")
	routes.RegisterTaskRoutes(taskGroup)

	addr := fmt.Sprintf(":%v", PORT)	
	log.Fatalln(r.Run(addr))
}