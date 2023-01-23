package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	DbContext struct {
		echo.Context
		Collection *mongo.Collection
	}
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

func checkErr(er error) {
	if er != nil {
		log.Fatal(er)
	}
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

func UpdateItem(collection *mongo.Collection, realintid int, c echo.Context) string {
	filterUpdate := bson.D{{Key: "id", Value: realintid}}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "successdetails", Value: c.FormValue("successdetails")},
			{Key: "title", Value: c.FormValue("title")},
			{Key: "details", Value: c.FormValue("details")},
			{Key: "successimage", Value: c.FormValue("successimage")},
			{Key: "category", Value: c.FormValue("category")},
			{Key: "success", Value: c.FormValue("success")},
			{Key: "date", Value: c.FormValue("date")},
			{Key: "image", Value: c.FormValue("image")},
		}},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filterUpdate, update)
	checkErr(err)

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	if updateResult.MatchedCount == 0 {
		return "ERR NO DOCUMENT(s)"
	}
	return "success"
}

func MakeItem(collection *mongo.Collection, c echo.Context) string {
	realintid, errs := strconv.Atoi(c.FormValue("id"))
	checkErr(errs)

	filterCheck := bson.D{{Key: "id", Value: realintid}}

	var result ListContent
	checkedErr := collection.FindOne(context.TODO(), filterCheck).Decode(&result)

	if checkedErr != nil {
		if checkedErr == mongo.ErrNoDocuments {
			ash := ListContent{c.FormValue("title"), c.FormValue("details"), c.FormValue("successdetails"), c.FormValue("successimage"), c.FormValue("category"), c.FormValue("success"), c.FormValue("date"), c.FormValue("image"), realintid}

			_, err := collection.InsertOne(context.TODO(), ash)
			checkErr(err)
			return "success"
		}
		panic(checkedErr)
	}
	return "DOCUMENT(s) ALEADY"
}
