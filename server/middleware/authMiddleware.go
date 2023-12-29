package middleware

import (
	"CRUD_API/server/models"
	utils "CRUD_API/server/utils"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	database "CRUD_API/server/settings"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var userSessionCollection *mongo.Collection = database.OpenCollection(database.Client, "userSession")

// func Authenticate() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		clientToken := c.Request.Header.Get("Authorization")
// 		if clientToken == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not provided"})
// 			c.Abort()
// 			return
// 		}

// 		claims, err := utils.ValidateToken(clientToken)
// 		if err != "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
// 			c.Abort()
// 			return
// 		}
// 		c.Set("email", claims.Email)
// 		c.Set("role", claims.Role)
// 		c.Set("firstName", claims.FirstName)
// 		// c.Set("uid", claims.Uid)
// 		c.Next()
// 	}
// }

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		fields := strings.Fields(clientToken)
		var access_token string
		if len(fields) != 0 && fields[0] == "Bearer" {
			access_token = fields[1]
		} else {
			access_token = ""
		}

		if access_token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token not provided"})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var userSession *models.UserSession

		sessionDataErr := userSessionCollection.FindOne(ctx, bson.M{"token": access_token}).Decode(&userSession)
		defer cancel()

		if sessionDataErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": sessionDataErr.Error()})
			return
		}

		var user *models.User
		userDataErr := userCollection.FindOne(ctx, bson.M{"_id": userSession.UserID}).Decode(&user)
		defer cancel()
		if userDataErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": userDataErr.Error()})
			return
		}

		// userSession.User = user

		claims, err := utils.ValidateToken(access_token)
		if err != "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}
		fmt.Println("claims", claims)
		// c.Set("email", claims.Email)
		// c.Set("role", claims.Role)
		// c.Set("firstName", claims.FirstName)
		// c.Set("userId", claims.UserId)
		c.Set("currentUser", user)
		c.Next()
	}
}
