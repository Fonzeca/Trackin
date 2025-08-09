package manager

import (
	"fmt"
	"time"

	"github.com/Fonzeca/Trackin/internal/core/services"
	"github.com/Fonzeca/Trackin/internal/infrastructure/database/model"
	"github.com/Fonzeca/Trackin/internal/infrastructure/database/query"
	"github.com/Fonzeca/Trackin/internal/infrastructure/geolocation"
	"github.com/Fonzeca/Trackin/internal/interfaces/messaging/json"
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
		fmt.Printf("[DATAENTRY] IMEI %s: Descartando datos con Latitude=0\n", simplyData.Imei)
		return nil
	}

	fmt.Printf("[DATAENTRY] IMEI %s: Procesando nuevo log - Date: %s, Lat: %f, Lng: %f, Speed: %f\n",
		simplyData.Imei, simplyData.Date.Add(-time.Hour*3).Format("2006-01-02 15:04:05"),
		simplyData.Latitude, simplyData.Longitude, simplyData.Speed)

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

	fmt.Printf("[DATAENTRY] IMEI %s: Guardando en DB - Date: %s\n",
		simplyData.Imei, log.Date.Format("2006-01-02 15:04:05"))

	err := q.Create(&log)
	if err != nil {
		fmt.Printf("[DATAENTRY] ERROR IMEI %s: Fallo al guardar en DB: %v\n", simplyData.Imei, err)
		return err
	}

	fmt.Printf("[DATAENTRY] IMEI %s: Log guardado exitosamente en DB\n", simplyData.Imei)

	lastPoint, ok := services.GetCachedPoints(simplyData.Imei)
	if !ok {
		fmt.Printf("[DATAENTRY] IMEI %s: No hay punto en cache, actualizando cache con nuevo log\n", simplyData.Imei)
		go services.SetCachedPoints(simplyData.Imei, &log)
	} else if lastPoint == nil {
		fmt.Printf("[DATAENTRY] IMEI %s: Punto en cache es nil, actualizando cache con nuevo log\n", simplyData.Imei)
		go services.SetCachedPoints(simplyData.Imei, &log)
	} else {
		fmt.Printf("[DATAENTRY] IMEI %s: Cache actual - Date: %s | Nuevo log - Date: %s\n",
			simplyData.Imei,
			lastPoint.Date.Format("2006-01-02 15:04:05"),
			log.Date.Format("2006-01-02 15:04:05"))

		isValid := geolocation.IsValidPoint(lastPoint, &log)
		fmt.Printf("[DATAENTRY] IMEI %s: IsValidPoint result: %t\n", simplyData.Imei, isValid)

		if isValid {
			fmt.Printf("[DATAENTRY] IMEI %s: Validación OK, actualizando cache\n", simplyData.Imei)
			go services.SetCachedPoints(simplyData.Imei, &log)
		} else {
			fmt.Printf("[DATAENTRY] IMEI %s: Validación FALLÓ, NO actualizando cache\n", simplyData.Imei)
		}
	}

	// d.geofenceService.DispatchMessage(simplyData)

	return nil
}
