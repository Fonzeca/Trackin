package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	manager "github.com/Fonzeca/Trackin/internal/core/managers"
	"github.com/Fonzeca/Trackin/internal/core/services"
	db "github.com/Fonzeca/Trackin/internal/infrastructure/database"
	model_json "github.com/Fonzeca/Trackin/internal/interfaces/messaging/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqDataEntry struct {
	trackinInputs <-chan amqp.Delivery
	rastroInputs  <-chan amqp.Delivery
}

func NewRabbitMqDataEntry() RabbitMqDataEntry {
	channel := services.GlobalRabbitChannel

	// Configurar cola 'trackin'
	qTrackin, err := channel.QueueDeclare("trackin", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = channel.QueueBind(qTrackin.Name, "trackin.data.*.decoded", "carmind", false, nil)
	if err != nil {
		panic(err)
	}

	// Configurar cola 'rastro'
	qRastro, err := channel.QueueDeclare("rastro", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	// Subscribing to queue 'trackin'
	trackinMessages, err := channel.Consume(
		qTrackin.Name, // queue name
		"trackin",     // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no local
		false,         // no wait
		nil,           // arguments
	)
	if err != nil {
		log.Println(err)
	}

	// Subscribing to queue 'rastro'
	rastroMessages, err := channel.Consume(
		qRastro.Name, // queue name
		"rastro",     // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no local
		false,        // no wait
		nil,          // arguments
	)
	if err != nil {
		log.Println(err)
	}

	instance := RabbitMqDataEntry{
		trackinInputs: trackinMessages,
		rastroInputs:  rastroMessages,
	}
	go instance.Run()

	return instance
}

func (m *RabbitMqDataEntry) Run() {
	entryManager := manager.GetManagerContainer().GetDataEntryManager()
	if entryManager == nil {
		fmt.Println("DataEntryManager is not initialized")
		return
	}

	// Usar select para escuchar ambas colas simultáneamente
	for {
		select {
		case message := <-m.trackinInputs:
			m.processTrackinMessage(message, entryManager)
		case message := <-m.rastroInputs:
			m.processRastroMessage(message, entryManager)
		}
	}
}

func (m *RabbitMqDataEntry) processTrackinMessage(message amqp.Delivery, entryManager manager.IDataEntryManager) {
	switch message.RoutingKey {
	case "trackin.data.log.decoded":
		pojo := model_json.SimplyData{}
		err := json.Unmarshal(message.Body, &pojo)
		if err != nil {
			fmt.Println("Error al deserializar el mensaje trackin:", string(message.Body))
			fmt.Println(err)
			message.Ack(false)
			return
		}

		pojo.PayLoad = string(message.Body)

		err = entryManager.ProcessData(pojo, db.DB)
		if err != nil {
			fmt.Println(err)
			return
		}
		message.Ack(false)
	default:
		fmt.Printf("Routing key no reconocido en trackin: %s\n", message.RoutingKey)
		message.Ack(false)
	}
}

func (m *RabbitMqDataEntry) processRastroMessage(message amqp.Delivery, entryManager manager.IDataEntryManager) {
	// Procesar mensajes de la cola 'rastro'
	fmt.Printf("Mensaje recibido de cola rastro - Routing Key: %s\n", message.RoutingKey)
	fmt.Printf("Contenido: %s\n", string(message.Body))

	pojo1 := model_json.SimplyDataLocation{}
	err := json.Unmarshal(message.Body, &pojo1)
	if err != nil {
		fmt.Println("Error al deserializar el mensaje rastro:", string(message.Body))
		fmt.Println(err)
		message.Ack(false)
		return
	}

	pojo := model_json.SimplyData{}
	pojo.Imei = pojo1.Imei
	pojo.Latitude = pojo1.Latitude
	pojo.Longitude = pojo1.Longitude
	pojo.Date = time.Unix(pojo1.TimestampMs/1000, (pojo1.TimestampMs%1000)*int64(time.Millisecond))
	pojo.Speed = float32(pojo1.Speed)
	pojo.EngineStatus = pojo1.EngineOn
	pojo.Azimuth = pojo1.Azimuth
	pojo.PayLoad = string(message.Body)

	err = entryManager.ProcessData(pojo, db.DB)
	if err != nil {
		fmt.Println(err)
		return
	}
	message.Ack(false)
}
