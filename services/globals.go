package services

import (
	"github.com/Fonzeca/Trackin/db/model"
	"github.com/rabbitmq/amqp091-go"
)

var (
	GlobalSender ISender

	GlobalRabbitChannel *amqp091.Channel

	CachedPoints map[string]*model.Log = make(map[string]*model.Log)
)
