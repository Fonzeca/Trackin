package entry

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Fonzeca/Trackin/db"
	model_json "github.com/Fonzeca/Trackin/entry/json"
	"github.com/Fonzeca/Trackin/services"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqDataEntry struct {
	inputs <-chan amqp.Delivery
}

func NewRabbitMqDataEntry() RabbitMqDataEntry {
	channel := services.GlobalRabbitChannel

	q, err := channel.QueueDeclare("trackin", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = channel.QueueBind(q.Name, "trackin.data.*.decoded", "carmind", false, nil)
	if err != nil {
		panic(err)
	}

	// Subscribing to QueueService1 for getting messages.
	messages, err := channel.Consume(
		q.Name,    // queue name
		"trackin", // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // arguments
	)
	if err != nil {
		log.Println(err)
	}

	instance := RabbitMqDataEntry{inputs: messages}
	go instance.Run()
	return instance
}

func (m *RabbitMqDataEntry) Run() {
	db, close, err := db.ObtenerConexionDb()
	if err != nil {
		log.Fatalln(err)
	}
	defer close()
	for message := range m.inputs {

		switch message.RoutingKey {
		case "trackin.data.log.decoded":
			pojo := model_json.SimplyData{}
			err := json.Unmarshal(message.Body, &pojo)
			if err != nil {
				fmt.Println("Error al deserializar el mensaje:", string(message.Body))
				fmt.Println(err)
				message.Ack(false)
				break
			}

			pojo.PayLoad = string(message.Body)

			err = DataEntryManager.ProcessData(pojo, db)
			if err != nil {
				fmt.Println(err)
				break
			}
			message.Ack(false)
			break
		}
	}
}
