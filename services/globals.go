package services

import "github.com/rabbitmq/amqp091-go"

var (
	GlobalSender ISender

	GlobalRabbitChannel *amqp091.Channel
)
