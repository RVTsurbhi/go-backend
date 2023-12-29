package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	models "CRUD_API/server/models"
	database "CRUD_API/server/settings"
	utils "CRUD_API/server/utils"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var userSessionCollection *mongo.Collection = database.OpenCollection(database.Client, "userSession")

func SignupUser(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	deviceType := c.Request.Header.Get("deviceType")
	deviceId := c.Request.Header.Get("deviceId")
	//NOTE: consider this to be body of the request
	var userInputData models.User

	if err := c.BindJSON(&userInputData); err != nil {
		e := utils.NewCustomError(http.StatusBadRequest, err.Error())
		c.Error(e)
		return
	}

	validate := validator.New()
	validationErr := validate.Struct(userInputData)
	if validationErr != nil {
		e := utils.NewCustomError(http.StatusBadRequest, validationErr.Error())
		c.Error(e)
		return
	}

	//NOTE: Email validation
	count, err := userCollection.CountDocuments(ctx, bson.M{"email": userInputData.Email})
	defer cancel()
	if err != nil {
		e := utils.NewCustomError(http.StatusBadGateway, err.Error())
		c.Error(e)
	}

	if count > 0 {
		e := utils.NewCustomError(http.StatusBadGateway, "Email Already Taken")
		c.Error(e)
		return
	}

	hashedPassword, _ := utils.HashPassword(userInputData.Password)
	userInputData.Password = hashedPassword
	userInputData.Role = "USER"
	userInputData.CreatedAt = time.Now()
	userInputData.UpdatedAt = time.Now()
	userInputData.ID = primitive.NewObjectID()

	_, insertErr := userCollection.InsertOne(ctx, userInputData)
	if insertErr != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"message": insertErr.Error()})
		e := utils.NewCustomError(http.StatusInternalServerError, insertErr.Error())
		c.Error(e)
		return
	}
	defer cancel()

	token, _ := utils.GenerateToken2(userInputData.Email, userInputData.FirstName, userInputData.Role, userInputData.ID.String())
	var sessionData models.UserSession
	sessionData.Token = token
	sessionData.UserID = userInputData.ID
	sessionData.DeviceID = deviceId
	sessionData.DeviceType = deviceType
	sessionData.ID = primitive.NewObjectID()

	//update user session
	_, updateErr := userSessionCollection.InsertOne(ctx, sessionData)
	defer cancel()

	if updateErr != nil {
		e := utils.NewCustomError(http.StatusInternalServerError, insertErr.Error())
		c.Error(e)
		return
	}
	utils.DataResponse(c, http.StatusOK, "User created successfully", gin.H{"data": token})
}

func SignInUser(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	//fetching headers from the request
	deviceType := c.Request.Header.Get("deviceType")
	deviceId := c.Request.Header.Get("deviceId")

	//NOTE: validate input request
	var credentials models.SignInInput
	var foundUser models.User

	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	err := userCollection.FindOne(ctx, bson.M{"email": credentials.Email}).Decode(&foundUser)
	defer cancel()
	if err != nil {
		e := utils.NewCustomError(http.StatusBadGateway, "Invalid email")
		c.Error(e)
		return
	}

	if err := utils.VerifyPassword(foundUser.Password, credentials.Password); err != nil {
		e := utils.NewCustomError(http.StatusBadGateway, "Invalid email or Password")
		c.Error(e)
		return
	}

	token, _ := utils.GenerateToken2(foundUser.Email, foundUser.FirstName, foundUser.Role, foundUser.ID.String())
	var sessionData models.UserSession
	sessionData.Token = token
	sessionData.UserID = foundUser.ID
	sessionData.DeviceID = deviceId
	sessionData.DeviceType = deviceType
	sessionData.ID = primitive.NewObjectID()

	//update user session
	_, updateErr := userSessionCollection.InsertOne(ctx, sessionData)
	defer cancel()

	if updateErr != nil {
		e := utils.NewCustomError(http.StatusInternalServerError, updateErr.Error())
		c.Error(e)
		return
	}

	utils.DataResponse(c, http.StatusOK, "User loggedIn successfully", gin.H{"data": token})
}

func GetUserProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := utils.MatchUserType(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": user})
	}
}
