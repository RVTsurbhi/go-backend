package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	models "CRUD_API/server/models"
	database "CRUD_API/server/settings"
	utils "CRUD_API/server/utils"
)

var postCollection *mongo.Collection = database.OpenCollection(database.Client, "post")

// create a post
func CreatePost(c *gin.Context) {
	//Note: retrieve the data set from middleware
	currentUser := c.MustGet("currentUser").(*models.User)

	fmt.Println("user", currentUser)
	var postPayload models.Post
	var newPost models.Post

	if err := c.BindJSON(&postPayload); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"message": err})
		// log.Fatal(err)

		e := utils.NewCustomError(http.StatusBadRequest, err.Error())
		c.Error(e)
		return
	}

	validate := validator.New()
	validationErr := validate.Struct(postPayload)
	if validationErr != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"message": validationErr.Error()})
		// return

		e := utils.NewCustomError(http.StatusBadRequest, validationErr.Error())
		c.Error(e)
		return
	}

	postPayload.CreatedAt = time.Now()
	postPayload.UpdatedAt = time.Now()
	postPayload.UserID = currentUser.ID

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	result, insertErr := postCollection.InsertOne(ctx, postPayload)
	defer cancel()
	if insertErr != nil {
		// log.Panic(insertErr)
		// c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while inserting"})
		e := utils.NewCustomError(http.StatusInternalServerError, "Error while inserting")
		c.Error(e)
		return
	}
	query := bson.M{"_id": result.InsertedID}
	foundPostErr := postCollection.FindOne(ctx, query).Decode(&newPost)
	defer cancel()
	if foundPostErr != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while fetching", "error": foundPostErr})
		e := utils.NewCustomError(http.StatusInternalServerError, foundPostErr.Error())
		c.Error(e)
		return
	}

	// c.JSON(http.StatusOK, gin.H{"message": "Post created", "data": newPost})
	utils.DataResponse(c, http.StatusOK, "Post created", gin.H{"data": newPost})
}

func GetPosts(c *gin.Context) {
	//NOTE: Get a query parameter with a default value
	var page = c.DefaultQuery("page", "1")
	var limit = c.DefaultQuery("limit", "10")

	intPage, err := strconv.Atoi(page)
	if err != nil {
		// c.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		e := utils.NewCustomError(http.StatusInternalServerError, err.Error())
		c.Error(e)
		return
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		// c.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		e := utils.NewCustomError(http.StatusInternalServerError, err.Error())
		c.Error(e)
		return
	}

	if intPage == 0 {
		intPage = 1
	}

	if intLimit == 0 {
		intLimit = 10
	}

	skip := (intPage - 1) * intLimit

	opt := options.FindOptions{}
	opt.SetLimit(int64(intLimit))
	opt.SetSkip(int64(skip))
	opt.SetSort(bson.M{"created_at": -1})

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	var posts []models.Post
	cursor, err := postCollection.Find(ctx, bson.M{}, &opt)
	defer cancel()
	if err != nil {
		// log.Panic(err)
		// c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while fetching"})
		e := utils.NewCustomError(http.StatusInternalServerError, "Error while fetching")
		c.Error(e)
		return
	}
	if err = cursor.All(ctx, &posts); err != nil {
		// log.Panic(err)
		// c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while fetching"})
		e := utils.NewCustomError(http.StatusInternalServerError, "Error while fetching")
		c.Error(e)
		return
	}

	count, countErr := postCollection.CountDocuments(ctx, bson.M{})
	defer cancel()
	if countErr != nil {
		// log.Panic(countErr)
		// c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while fetching"})
		e := utils.NewCustomError(http.StatusInternalServerError, countErr.Error())
		c.Error(e)
		return
	}
	// c.JSON(http.StatusOK, gin.H{"data": posts, "totalCount": count})
	utils.PageResponse(c, http.StatusOK, count, skip, gin.H{"data": posts})
}

func GetPostById(c *gin.Context) {
	postId := c.Param("postId")
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	var post models.Post
	obId, _ := primitive.ObjectIDFromHex(postId)
	query := bson.M{"_id": obId}
	err := postCollection.FindOne(ctx, query).Decode(&post)
	defer cancel()
	if err != nil {
		// log.Panic(err)
		// c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while fetching"})
		e := utils.NewCustomError(http.StatusInternalServerError, "Error while fetching")
		c.Error(e)
		return
	}
	// c.JSON(http.StatusOK, gin.H{"data": post})
	utils.DataResponse(c, http.StatusOK, "", gin.H{"data": post})
}

func UpdatePost(c *gin.Context) {
	//Note: retrieve the data set from middleware
	currentUser := c.MustGet("currentUser").(*models.User)
	postId := c.Param("postId")
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	var postPayload models.Post
	var updatedPost models.Post
	if err := c.BindJSON(&postPayload); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"message": err})
		// log.Fatal(err)

		e := utils.NewCustomError(http.StatusBadRequest, err.Error())
		c.Error(e)
		return
	}
	validate := validator.New()
	validationErr := validate.Struct(postPayload)
	if validationErr != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"message": validationErr.Error()})
		e := utils.NewCustomError(http.StatusBadRequest, validationErr.Error())
		c.Error(e)
		return
	}
	obId, _ := primitive.ObjectIDFromHex(postId)
	query := bson.M{"_id": obId, "userId": currentUser.ID}
	update := bson.M{
		"$set": bson.M{
			"title":      postPayload.Title,
			"content":    postPayload.Content,
			"updated_at": time.Now(),
		},
	}
	err := postCollection.FindOneAndUpdate(ctx, query, update).Decode(&updatedPost)
	defer cancel()
	if err != nil {
		// log.Panic(err)
		// c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while updating"})
		e := utils.NewCustomError(http.StatusInternalServerError, "Error while updating")
		c.Error(e)
		return
	}
	// c.JSON(http.StatusOK, gin.H{"message": "Post updated", "data": updatedPost})
	utils.DataResponse(c, http.StatusOK, "Post updated", gin.H{"data": updatedPost})
}
