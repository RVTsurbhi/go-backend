package main

import (
	"log"
	"net/http"

	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	routes "CRUD_API/server/routes"
	utils "CRUD_API/server/utils"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	server := gin.Default()
	server.Use(gin.Logger())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "https://another.example.com"}

	server.Use(cors.New(corsConfig))
	server.Use(utils.ErrorHandler())

	server.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "server is running"})
	})

	server.NoRoute(func(c *gin.Context) {
		err := utils.NewCustomError(http.StatusNotFound, "Not Found")
		c.Error(err)
	})

	routes.AuthRoutes(server)
	routes.UserRoutes(server)
	routes.PostRoutes(server)

	server.Run(":" + port)
}
