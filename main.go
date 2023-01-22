package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

var mongoSettings string = "몽글몽글한 몽고"

func main() {
	clientOptions := options.Client().ApplyURI(mongoSettings)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	Collection := client.Database("bucketlist").Collection("items")

	/*content, err := ioutil.ReadFile("./d.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload []ListContent
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	fmt.Println(payload[0].Title)

	for _, value := range payload {
		fmt.Println(value)
		insertResult, err := Collection.InsertOne(context.TODO(), value)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	} */

	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/updateitemsuccessdetails", func(c echo.Context) error {
		passkey := c.FormValue("passcodeforme")
		if passkey != "패스하던가 말던가 ㅋ" {
			return c.String(http.StatusUnauthorized, "No Permission")
		}
		realintid, _ := strconv.Atoi(c.FormValue("id"))

		filterUpdate := bson.D{{Key: "id", Value: realintid}}

		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "successdetails", Value: c.FormValue("successdetails")},
			}},
		}
		updateResult, err := Collection.UpdateOne(context.TODO(), filterUpdate, update)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
		return c.String(http.StatusOK, "successful")
	})
	e.GET("/getall", func(c echo.Context) error {
		findOptions := options.Find()

		var results []*ListContent

		cur, err := Collection.Find(context.TODO(), bson.D{{}}, findOptions)
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

		return c.String(http.StatusOK, string(out))
	})
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "https://bucket.yeonw.me"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	e.Logger.Fatal(e.Start(":1323"))
}
