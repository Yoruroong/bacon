package main

import (
	"net/http"

	"github.com/Yoruroong/bacon/db"
	"github.com/Yoruroong/bacon/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var mongoSettings string = "몽글몽글한 몽고"

func main() {
	collection := db.Connect(mongoSettings)

	e := echo.New()
	e.GET("/", routes.HandleHome)
	e.GET("/getall", routes.GetAll)
	e.POST("/updateitem", routes.UpdateItem)
	e.POST("/makeitem", routes.MakeItem)

	e.Use(middleware.Logger())
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &db.DbContext{Context: c, Collection: collection}
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
