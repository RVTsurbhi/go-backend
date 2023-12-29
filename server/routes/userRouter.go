package routes

import (
	controller "CRUD_API/server/controllers"
	middleware "CRUD_API/server/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(rc *gin.Engine) {
	router := rc.Group("/api/user")

	router.Use(middleware.Authenticate())

	router.GET("/profile/:user_id", controller.GetUserProfile())
}
