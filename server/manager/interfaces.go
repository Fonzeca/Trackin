package manager

import (
	"github.com/Fonzeca/Trackin/db/model"
)

// IRoutesManager define la interfaz para el manager de rutas
type IRoutesManager interface {
	GetLastLogByImei(imei string) (model.LastLogView, error)
	GetVehiclesStateByImeis(only string, imeis model.ImeisBody) ([]model.StateLogView, error)
	GetRouteByImei(requestRoute model.RouteRequest) ([]model.GpsRouteData, error)
	GetRouteByImeiAndZones(requestRoute model.RouteRequest, zones []model.ZoneView) ([]model.GpsRouteData, error)
	CleanUpRouteBySpeedAnomaly(route []model.GpsPoint) []model.GpsPoint

	// Setter para inyección de dependencias
	SetZonasManager(zonasManager IZonasManager)
}

// IDataEntryManager define la interfaz para el manager de entrada de datos
type IDataEntryManager interface {
	ProcessData(data interface{}, db interface{}) error

	// Setters para inyección de dependencias
	SetRoutesManager(routesManager IRoutesManager)
	SetZonasManager(zonasManager IZonasManager)
}
