package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	uri := os.Getenv("DB_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	defer cancel()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to mongoDB")
	return client
}

var Client *mongo.Client = ConnectDB()
var DBinstance *mongo.Database = nil

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	if DBinstance == nil {
		dbName := os.Getenv("DB_NAME")
		DBinstance = client.Database(dbName)
	}

	var collection *mongo.Collection = DBinstance.Collection(collectionName)
	return collection
}
