package main

import (
	"fmt"

	"github.com/Fonzeca/Trackin/entry"
	"github.com/Fonzeca/Trackin/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	InitConfig()
	channel, closeFunc := setupRabbitMq()
	defer closeFunc()
	entry.NewRabbitMqDataEntry(channel)

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

func setupRabbitMq() (*amqp.Channel, func()) {
	// Create a new RabbitMQ connection.

	connectRabbitMQ, err := amqp.Dial(viper.GetString("rabbitmq.url"))
	if err != nil {
		panic(err)
	}

	// Opening a channel to our RabbitMQ instance over
	// the connection we have already established.
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		connectRabbitMQ.Close()
		panic(err)
	}

	return channelRabbitMQ, func() { connectRabbitMQ.Close(); channelRabbitMQ.Close() }
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
