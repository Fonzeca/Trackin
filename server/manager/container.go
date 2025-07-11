package manager

import (
	"sync"
)

// ManagerContainer es el contenedor singleton que maneja todas las instancias de managers
type ManagerContainer struct {
	routesManager    IRoutesManager
	zonasManager     IZonasManager
	dataEntryManager IDataEntryManager
	initialized      bool
	mu               sync.RWMutex
}

var (
	instance *ManagerContainer
	once     sync.Once
)

// GetManagerContainer devuelve la instancia singleton del contenedor
func GetManagerContainer() *ManagerContainer {
	once.Do(func() {
		instance = &ManagerContainer{}
		instance.initialize()
	})
	return instance
}

// initialize inicializa todas las dependencias y las conecta
func (c *ManagerContainer) initialize() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return
	}

	// Crear las instancias
	c.routesManager = newRoutesManager()
	c.zonasManager = newZonasManager()
	c.dataEntryManager = newDataEntryManager()

	// Inyectar dependencias
	c.routesManager.SetZonasManager(c.zonasManager)
	c.zonasManager.SetRoutesManager(c.routesManager)

	c.dataEntryManager.SetRoutesManager(c.routesManager)
	c.dataEntryManager.SetZonasManager(c.zonasManager)

	c.initialized = true
}

// GetRoutesManager devuelve la instancia del routes manager
func (c *ManagerContainer) GetRoutesManager() IRoutesManager {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.routesManager
}

// GetZonasManager devuelve la instancia del zonas manager
func (c *ManagerContainer) GetZonasManager() IZonasManager {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.zonasManager
}

// GetDataEntryManager devuelve la instancia del data entry manager
func (c *ManagerContainer) GetDataEntryManager() IDataEntryManager {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.dataEntryManager
}

// Reset reinicia el contenedor (útil para testing)
func (c *ManagerContainer) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.routesManager = nil
	c.zonasManager = nil
	c.dataEntryManager = nil
	c.initialized = false
}

// ResetSingleton reinicia completamente el singleton (útil para testing)
func ResetSingleton() {
	instance = nil
	once = sync.Once{}
}
