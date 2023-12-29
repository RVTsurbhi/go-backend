package routes

import (
	controller "CRUD_API/server/controllers"
	middleware "CRUD_API/server/middleware"

	"github.com/gin-gonic/gin"
)

func PostRoutes(rc *gin.Engine) {
	router := rc.Group("/api/posts")

	// router.Use(middleware.ErrorHandler())
	router.Use(middleware.Authenticate())

	router.POST("/", controller.CreatePost)
	router.GET("/list", controller.GetPosts)
	router.GET("/:postId", controller.GetPostById)
	router.PATCH("/:postId", controller.UpdatePost)
	// router.DELETE("/:postId", r.postController.DeletePost)
}
