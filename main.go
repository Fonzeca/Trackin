package main

import (
	"net/http"

	"github.com/Fonzeca/Trackin/db"
	"github.com/Fonzeca/Trackin/rest/json"
	"github.com/labstack/echo/v4"
)

func main() {

	canal := make(chan json.SimplyData)

	go db.Deamon(canal)

	e := echo.New()
	e.POST("/data", func(c echo.Context) error {
		data := json.SimplyData{}
		c.Bind(&data)
		canal <- data
		return c.NoContent(http.StatusOK)
	})
	e.Logger.Fatal(e.Start(":4762"))
}
