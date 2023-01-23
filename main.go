package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Yoruroong/bacon/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

type (
	dbContext struct {
		echo.Context
		collection *mongo.Collection
	}
)

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func getAll(c echo.Context) error {
	collection := c.(*dbContext).collection
	items := db.GetAll(collection)

	return c.String(http.StatusOK, items)
}

var mongoSettings string = "몽글몽글한 몽고"

func main() {
	collection := db.Connect(mongoSettings)

	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/updateitemsuccessdetails", func(c echo.Context) error {
		passkey := c.FormValue("passcodeforme")
		if passkey != "과연 알 수 있을까?" {
			return c.String(http.StatusUnauthorized, "No Permission")
		}
		realintid, _ := strconv.Atoi(c.FormValue("id"))

		filterUpdate := bson.D{{Key: "id", Value: realintid}}

		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "successdetails", Value: c.FormValue("successdetails")},
			}},
		}
		updateResult, err := collection.UpdateOne(context.TODO(), filterUpdate, update)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
		return c.String(http.StatusOK, "successful")
	})
	e.GET("/getall", getAll)

	e.Use(middleware.Logger())
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &dbContext{c, collection}
			return h(cc)
		}
	})
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "https://bucket.yeonw.me"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	e.Logger.Fatal(e.Start(":1323"))
}
