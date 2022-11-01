package services

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

func SetupRabbitMq() (*amqp.Channel, func()) {
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
	channelRabbitMQ.Qos(10, 0, false)

	GlobalSender = &SenderRabbitMq{
		channel: channelRabbitMQ,
	}
	GlobalRabbitChannel = channelRabbitMQ
	return channelRabbitMQ, func() { connectRabbitMQ.Close(); channelRabbitMQ.Close() }
}

type SenderRabbitMq struct {
	channel *amqp.Channel
}

func (s *SenderRabbitMq) SendMessage(context context.Context, destination string, message []byte) error {
	return s.channel.PublishWithContext(context, "carmind", destination, false, false, amqp.Publishing{
		Body:        message,
		ContentType: "application/json",
	})
}
