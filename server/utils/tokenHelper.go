package utils

import (
	"context"
	"log"
	"os"
	"time"

	database "CRUD_API/server/settings"

	jwt "github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email     string
	FirstName string
	Role      string
	UserId    string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var JWT_SECRET_KEY string = os.Getenv("JWT_SECRET_KEY")

func GenerateToken2(email string, firstName string, role string, userId string) (string, error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		Role:      role,
		UserId:    userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(JWT_SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return "", err
	}
	return token, err
}

func ValidateToken(token string) (claims *SignedDetails, msg string) {
	claims = &SignedDetails{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		msg = "Token expired"
		return
	}
	return claims, msg
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	var context, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refreshToken", signedRefreshToken})

	Updated_at := time.Now()
	updateObj = append(updateObj, bson.E{"updated_at", Updated_at})

	upsert := true
	filter := bson.M{"_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(context, filter, bson.D{{"$set", updateObj}}, &opt)

	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
}
