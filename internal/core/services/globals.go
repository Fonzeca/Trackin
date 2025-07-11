package services

import (
	"sync"

	"github.com/Fonzeca/Trackin/internal/core/domain"
	"github.com/Fonzeca/Trackin/internal/infrastructure/database/model"
	"github.com/rabbitmq/amqp091-go"
)

var (
	GlobalSender domain.ISender

	GlobalRabbitChannel *amqp091.Channel

	cachedPoints      map[string]*model.Log = make(map[string]*model.Log)
	cachedPointsMutex                       = sync.RWMutex{}

	// Locks por IMEI para evitar consultas duplicadas
	imeiLocks      map[string]*sync.Mutex = make(map[string]*sync.Mutex)
	imeiLocksMutex                        = sync.Mutex{}
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

// GetOrCreateImeiLock obtiene o crea un lock específico para un IMEI
func GetOrCreateImeiLock(imei string) *sync.Mutex {
	imeiLocksMutex.Lock()
	defer imeiLocksMutex.Unlock()

	if lock, exists := imeiLocks[imei]; exists {
		return lock
	}

	lock := &sync.Mutex{}
	imeiLocks[imei] = lock
	return lock
}

// GetCachedPointsWithLock obtiene puntos del caché con lock para evitar consultas duplicadas
func GetCachedPointsWithLock(imei string, dbQueryFunc func() *model.Log) (*model.Log, bool) {
	// Primero verificamos si ya está en caché sin lock
	if point, ok := GetCachedPoints(imei); ok && point != nil {
		return point, true
	}

	// Si no está en caché, obtenemos el lock específico para este IMEI
	lock := GetOrCreateImeiLock(imei)
	lock.Lock()
	defer lock.Unlock()

	// Verificamos nuevamente por si otra goroutine ya lo cargó
	if point, ok := GetCachedPoints(imei); ok && point != nil {
		return point, true
	}

	// Si aún no está, ejecutamos la consulta a la DB
	point := dbQueryFunc()
	if point != nil {
		SetCachedPoints(imei, point)
	}

	return point, false
}
