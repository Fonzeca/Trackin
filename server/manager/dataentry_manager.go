package manager

import (
	"time"

	"github.com/Fonzeca/Trackin/db/model"
	"github.com/Fonzeca/Trackin/db/query"
	"github.com/Fonzeca/Trackin/entry/json"
	"github.com/Fonzeca/Trackin/services"
	"gorm.io/gorm"
)

type dataEntryManager struct {
	routesManager IRoutesManager
	zonasManager  IZonasManager
}

func newDataEntryManager() IDataEntryManager {
	return &dataEntryManager{}
}

// SetRoutesManager inyecta la dependencia del routes manager
func (d *dataEntryManager) SetRoutesManager(routesManager IRoutesManager) {
	d.routesManager = routesManager
}

// SetZonasManager inyecta la dependencia del zonas manager
func (d *dataEntryManager) SetZonasManager(zonasManager IZonasManager) {
	d.zonasManager = zonasManager
}

func (d *dataEntryManager) ProcessData(data interface{}, db interface{}) error {
	simplyData, ok := data.(json.SimplyData)
	if !ok {
		return nil // Si no es el tipo esperado, no hacer nada
	}

	gormDB, ok := db.(*gorm.DB)
	if !ok {
		return nil // Si no es el tipo esperado, no hacer nada
	}

	//Evitamos datos inecesarios que llegan por equivocacion.
	if simplyData.Latitude == 0 {
		return nil
	}

	log := model.Log{
		Imei:         simplyData.Imei,
		ProtocolType: simplyData.ProtocolType,
		Latitud:      simplyData.Latitude,
		Longitud:     simplyData.Longitude,
		Date:         simplyData.Date.Add(-time.Hour * 3),
		Speed:        simplyData.Speed,
		AnalogInput1: simplyData.AnalogInput1,
		DeviceTemp:   simplyData.DeviceTemp,
		Mileage:      simplyData.Mileage,
		IsGps:        simplyData.GpsWorking,
		IsHistory:    simplyData.IsHistory,
		EngineStatus: simplyData.EngineStatus,
		Azimuth:      simplyData.Azimuth,
		Payload:      simplyData.PayLoad,
	}

	q := query.Use(gormDB).Log

	err := q.Create(&log)
	if err != nil {
		return err
	}

	lastPoint, ok := services.GetCachedPoints(simplyData.Imei)
	if ok && services.IsValidPoint(lastPoint, &log) {
		go services.SetCachedPoints(simplyData.Imei, &log)
	}

	// d.geofenceService.DispatchMessage(simplyData)

	return nil
}
