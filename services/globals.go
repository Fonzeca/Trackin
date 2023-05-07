package services

import (
	"sync"

	"github.com/Fonzeca/Trackin/db/model"
	"github.com/rabbitmq/amqp091-go"
)

var (
	GlobalSender ISender

	GlobalRabbitChannel *amqp091.Channel

	cachedPoints      map[string]*model.Log = make(map[string]*model.Log)
	cachedPointsMutex                       = sync.RWMutex{}
)

func GetCachedPoints(key string) (*model.Log, bool) {
	cachedPointsMutex.RLock()
	defer cachedPointsMutex.RUnlock()
	value, ok := cachedPoints[key]
	return value, ok
}

func SetCachedPoints(key string, value *model.Log) {
	cachedPointsMutex.Lock()
	defer cachedPointsMutex.Unlock()

	lastpoint, ok := cachedPoints[key]
	if !ok || lastpoint == nil {
		cachedPoints[key] = value
		return
	}
	if lastpoint.Date.Before(value.Date) {
		cachedPoints[key] = value
	}
}
