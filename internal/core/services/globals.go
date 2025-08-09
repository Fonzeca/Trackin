package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/Fonzeca/Trackin/internal/core/domain"
	"github.com/Fonzeca/Trackin/internal/infrastructure/database/model"
	"github.com/rabbitmq/amqp091-go"
)

// Tipos para el sistema de caché basado en canales
type cacheOperation string

const (
	opGet        cacheOperation = "get"
	opSet        cacheOperation = "set"
	opGetBatch   cacheOperation = "getBatch"
	opGetOrQuery cacheOperation = "getOrQuery"
)

type cacheRequest struct {
	operation cacheOperation
	key       string
	keys      []string // Para operaciones batch
	value     *model.Log
	queryFunc func() *model.Log // Para operaciones getOrQuery
	response  chan cacheResponse
}

type cacheResponse struct {
	value  *model.Log
	values map[string]*model.Log // Para operaciones batch
	found  bool
}

type CacheManager struct {
	requests chan cacheRequest
	cache    map[string]*model.Log
	done     chan struct{}
}

var (
	GlobalSender domain.ISender

	GlobalRabbitChannel *amqp091.Channel

	// Cache manager singleton
	cacheManager *CacheManager
	initOnce     sync.Once
)

// InitCacheManager inicializa el cache manager singleton
func InitCacheManager() *CacheManager {
	initOnce.Do(func() {
		cacheManager = &CacheManager{
			requests: make(chan cacheRequest, 100), // Buffer para reducir bloqueos
			cache:    make(map[string]*model.Log),
			done:     make(chan struct{}),
		}
		go cacheManager.worker()
	})
	return cacheManager
}

// GetCacheManager obtiene la instancia del cache manager
func GetCacheManager() *CacheManager {
	if cacheManager == nil {
		return InitCacheManager()
	}
	return cacheManager
}

// worker procesa todas las operaciones del caché en una sola goroutine
func (cm *CacheManager) worker() {
	for {
		select {
		case req := <-cm.requests:
			switch req.operation {
			case opGet:
				value, exists := cm.cache[req.key]
				req.response <- cacheResponse{
					value: value,
					found: exists,
				}
			case opSet:
				// Solo actualizar si es más reciente o no existe
				if existing, exists := cm.cache[req.key]; !exists || existing == nil {
					cm.cache[req.key] = req.value
				} else if req.value != nil && existing.Date.Before(req.value.Date) {
					cm.cache[req.key] = req.value
				}
				req.response <- cacheResponse{found: true}
			case opGetBatch:
				values := make(map[string]*model.Log)
				for _, key := range req.keys {
					if value, exists := cm.cache[key]; exists {
						values[key] = value
					}
				}
				req.response <- cacheResponse{
					values: values,
					found:  len(values) > 0,
				}
			case opGetOrQuery:
				// Verificar primero si está en caché
				if value, exists := cm.cache[req.key]; exists && value != nil {
					req.response <- cacheResponse{
						value: value,
						found: true,
					}
				} else {
					// No está en caché, ejecutar la consulta
					result := req.queryFunc()
					if result != nil {
						// Guardar en caché
						cm.cache[req.key] = result
					}
					req.response <- cacheResponse{
						value: result,
						found: result != nil,
					}
				}
			}
		case <-cm.done:
			return
		}
	}
}

// GetCachedPoints obtiene un punto del caché usando canales
func GetCachedPoints(key string) (*model.Log, bool) {
	cm := GetCacheManager()
	response := make(chan cacheResponse, 1)

	select {
	case cm.requests <- cacheRequest{
		operation: opGet,
		key:       key,
		response:  response,
	}:
		select {
		case resp := <-response:
			return resp.value, resp.found
		case <-time.After(5 * time.Second):
			fmt.Printf("Timeout getting cached point for key: %s\n", key)
			return nil, false
		}
	case <-time.After(time.Second):
		fmt.Printf("Cache manager busy, timeout enqueueing get request for key: %s\n", key)
		return nil, false
	}
}

// SetCachedPoints establece un punto en el caché usando canales
func SetCachedPoints(key string, value *model.Log) {
	cm := GetCacheManager()
	response := make(chan cacheResponse, 1)

	select {
	case cm.requests <- cacheRequest{
		operation: opSet,
		key:       key,
		value:     value,
		response:  response,
	}:
		select {
		case <-response:
			return
		case <-time.After(5 * time.Second):
			fmt.Printf("Timeout setting cached point for key: %s\n", key)
			return
		}
	case <-time.After(time.Second * 5):
		fmt.Printf("Cache manager busy, timeout 5 sec enqueueing set request for key: %s\n", key)
		return
	}
}

// StopCacheManager detiene el cache manager de forma segura
func StopCacheManager() {
	if cacheManager != nil {
		close(cacheManager.done)
	}
}

// GetCachedPointsBatch obtiene múltiples puntos del caché de una vez
func GetCachedPointsBatch(keys []string) (map[string]*model.Log, bool) {
	if len(keys) == 0 {
		return make(map[string]*model.Log), false
	}

	cm := GetCacheManager()
	response := make(chan cacheResponse, 1)

	select {
	case cm.requests <- cacheRequest{
		operation: opGetBatch,
		keys:      keys,
		response:  response,
	}:
		select {
		case resp := <-response:
			return resp.values, resp.found
		case <-time.After(5 * time.Second):
			fmt.Printf("Timeout getting cached points batch for %d keys\n", len(keys))
			return make(map[string]*model.Log), false
		}
	case <-time.After(time.Second):
		fmt.Printf("Cache manager busy, timeout enqueueing batch get request for %d keys\n", len(keys))
		return make(map[string]*model.Log), false
	}
}

// GetCachedPointsWithQuery obtiene puntos del caché o ejecuta consulta usando el cache manager
func GetCachedPointsWithQuery(imei string, dbQueryFunc func() *model.Log) (*model.Log, bool) {
	cm := GetCacheManager()
	response := make(chan cacheResponse, 1)

	select {
	case cm.requests <- cacheRequest{
		operation: opGetOrQuery,
		key:       imei,
		queryFunc: dbQueryFunc,
		response:  response,
	}:
		select {
		case resp := <-response:
			return resp.value, resp.found
		case <-time.After(10 * time.Second):
			fmt.Printf("Timeout getting or querying cached point for key: %s\n", imei)
			return nil, false
		}
	case <-time.After(time.Second):
		fmt.Printf("Cache manager busy, timeout enqueueing getOrQuery request for key: %s\n", imei)
		return nil, false
	}
}
