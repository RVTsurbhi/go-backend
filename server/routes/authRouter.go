package routes

import (
	controller "CRUD_API/server/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(rc *gin.Engine) {
	router := rc.Group("/api")

	router.POST("/signup", controller.SignupUser)
	router.POST("/login", controller.SignInUser)
	// rc.GET("/refresh", authController.RefreshAccessToken)
	// rc.GET("/logout", middleware.DeserializeUser(userService), rc.authController.LogoutUser)
	// rc.GET("/verifyemail/:verificationCode", authController.VerifyEmail)
	// rc.POST("/forgotPassword", authController.ForgotPassword)
}
