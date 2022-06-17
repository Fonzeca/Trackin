package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Fonzeca/Trackin/db"
	jsonModel "github.com/Fonzeca/Trackin/rest/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {

	canal := make(chan jsonModel.SimplyData)

	go db.Deamon(canal)

	e := echo.New()
	e.POST("/data", func(c echo.Context) error {
		data := jsonModel.SimplyData{}
		list := GetJSONRawBody(c)
		for k, _ := range list {
			fmt.Println(k)
		}

		c.Bind(&data)
		canal <- data
		return c.NoContent(http.StatusOK)
	})
	e.Logger.Fatal(e.Start(":4762"))
}

func GetJSONRawBody(c echo.Context) map[string]interface{} {

	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)

	if err != nil {

		log.Error("empty json body")
		return nil
	}

	return jsonBody
}
