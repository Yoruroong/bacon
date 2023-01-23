package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Yoruroong/bacon/db"
	"github.com/labstack/echo/v4"
)

func checkErr(er error) {
	if er != nil {
		log.Fatal(er)
	}
}

func HandleHome(c echo.Context) error {
	return c.File("home.html")
}

func GetAll(c echo.Context) error {
	collection := c.(*db.DbContext).Collection
	items := db.GetAll(collection)

	return c.String(http.StatusOK, items)
}

func UpdateItem(c echo.Context) error {
	collection := c.(*db.DbContext).Collection
	passkey := c.FormValue("passcodeforme")
	if passkey != "youshallnotpass" {
		return c.String(http.StatusUnauthorized, "No Permission")
	}
	if len(c.FormValue("id")) < 1 {
		return c.String(http.StatusTeapot, "Insert ID")
	}
	realintid, errs := strconv.Atoi(c.FormValue("id"))
	checkErr(errs)

	resultstr := db.UpdateItem(collection, realintid, c)

	return c.String(http.StatusOK, resultstr)
}

func MakeItem(c echo.Context) error {
	collection := c.(*db.DbContext).Collection
	passkey := c.FormValue("passcodeforme")
	if passkey != "youshallnotpass" {
		return c.String(http.StatusUnauthorized, "No Permission")
	}
	if len(c.FormValue("id")) < 1 {
		return c.String(http.StatusTeapot, "Insert ID")
	}
	resultstr := db.MakeItem(collection, c)

	return c.String(http.StatusOK, resultstr)
}
