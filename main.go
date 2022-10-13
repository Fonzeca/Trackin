package main

import (
	"fmt"

	"github.com/Fonzeca/Trackin/entry"
	"github.com/Fonzeca/Trackin/server"
	"github.com/Fonzeca/Trackin/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

func main() {
	InitConfig()
	_, closeFunc := services.SetupRabbitMq()
	defer closeFunc()
	entry.NewRabbitMqDataEntry()

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

func InitConfig() {
	viper.SetConfigName("config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}
