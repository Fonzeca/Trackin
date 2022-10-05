package main

import (
	"encoding/json"

	"github.com/Fonzeca/Trackin/entry"
	"github.com/Fonzeca/Trackin/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	entry.Router(e)

	api := server.NewApi()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	e.GET("/getLastLogByImei", api.GetLastLogByImei)
	e.POST("/getVehiclesStateByImeis", api.GetVehiclesStateByImeis)
	e.POST("/getRouteByImei", api.GetRouteByImei)

	e.GET("/getZonesByEmpresaId", api.GetZonesByEmpresaId)
	e.POST("/createZone", api.CreateZone)
	e.PUT("/editZoneById", api.EditZoneById)
	e.DELETE("/deleteZoneById", api.DeleteZoneById)

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
