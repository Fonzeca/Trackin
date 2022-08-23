package main

import (
	"encoding/json"

	"github.com/Fonzeca/Trackin/entry"
	"github.com/Fonzeca/Trackin/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	entry.Router(e)

	api := server.NewApi()

	e.GET("/getLastLogByImei", api.GetLastLogByImei)
	e.GET("/getVehiclesStateByImeis", api.GetVehiclesStateByImeis)
	e.GET("/getRouteByImei", api.GetRouteByImei)

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
