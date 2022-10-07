package entry

import (
	"encoding/json"
	"fmt"
	"log"

	model_json "github.com/Fonzeca/Trackin/entry/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqDataEntry struct {
	inputs <-chan amqp.Delivery
}

func NewRabbitMqDataEntry(channel *amqp.Channel) RabbitMqDataEntry {

	q, err := channel.QueueDeclare("trackin", false, true, false, false, nil)
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
	for message := range m.inputs {

		switch message.RoutingKey {
		case "trackin.data.log.decoded":
			pojo := model_json.SimplyData{}
			err := json.Unmarshal(message.Body, &pojo)
			if err != nil {
				fmt.Println(err)
				break
			}

			err = DataEntryManager.ProcessData(pojo)
			if err != nil {
				fmt.Println(err)
				break
			}

			message.Ack(false)
			break
		}
	}
}
