package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionstring = "mongodb+srv://immigrantsadmin:dipankar@cluster0.zpi714a.mongodb.net/?retryWrites=true&w=majority"

func DBinstance() *mongo.Client {

	clientOptions := options.Client().ApplyURI(connectionstring)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("ims").Collection(collectionName)
	return collection
}

func FeedbackCollection() *mongo.Collection {
	collection := Client.Database("ims").Collection("feedback")
	return collection
}
