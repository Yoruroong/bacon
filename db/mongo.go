package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListContent struct {
	// MongoID        string
	Title          string
	Details        string
	Successdetails string
	Successimage   string
	Category       string
	Success        string
	Date           string
	Image          string
	Id             int
}

func Connect(mongoSetting string) *mongo.Collection {
	clientOptions := options.Client().ApplyURI(mongoSetting)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	collection := client.Database("bucketlist").Collection("items")
	return collection
}

func GetAll(collection *mongo.Collection) string {
	findOptions := options.Find()

	var results []*ListContent

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem ListContent
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	out, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}
	return string(out)
}
