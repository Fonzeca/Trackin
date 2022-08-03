package main

import (
	"encoding/json"

	"github.com/Fonzeca/Trackin/entry"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	entry.Router(e)

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
