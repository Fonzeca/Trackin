package manager

import (
	"github.com/Fonzeca/Trackin/db/model"
	"github.com/Fonzeca/Trackin/db/query"
	"github.com/Fonzeca/Trackin/entry/json"
	"gorm.io/gorm"
)

type DataEntryManager struct {
	geofenceService GeofenceDetector
}

func NewDataEntryManager() *DataEntryManager {
	instance := &DataEntryManager{
		geofenceService: *NewGeofenceDetector(),
	}

	return instance
}

func (d *DataEntryManager) ProcessData(data json.SimplyData, db *gorm.DB) error {
	//Evitamos datos inecesarios que llegan por equivocacion.
	if data.Latitude == 0 {
		return nil
	}

	log := model.Log{
		Imei:         data.Imei,
		ProtocolType: data.ProtocolType,
		Latitud:      data.Latitude,
		Longitud:     data.Longitude,
		Date:         data.Date,
		Speed:        data.Speed,
		AnalogInput1: data.AnalogInput1,
		DeviceTemp:   data.DeviceTemp,
		Mileage:      data.Mileage,
		IsGps:        data.GpsWorking,
		IsHistory:    data.IsHistory,
		EngineStatus: data.EngineStatus,
		Azimuth:      data.Azimuth,
		Payload:      data.PayLoad,
	}

	q := query.Use(db).Log

	err := q.Create(&log)
	if err != nil {
		return err
	}

	d.geofenceService.DispatchMessage(data)

	return nil
}
