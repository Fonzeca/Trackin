package manager

import (
	"time"

	"github.com/Fonzeca/Trackin/db/model"
	"github.com/Fonzeca/Trackin/db/query"
	"github.com/Fonzeca/Trackin/entry/json"
	"gorm.io/gorm"
)

var argTimeZone *time.Location

type DataEntryManager struct {
	geofenceService GeofenceDetector
}

func setTimeZone() {
	arg, errArg := time.LoadLocation("America/Argentina/Buenos_Aires")
	if errArg != nil {
		panic(errArg)
	}
	argTimeZone = arg
}

func NewDataEntryManager() *DataEntryManager {
	instance := &DataEntryManager{
		geofenceService: *NewGeofenceDetector(),
	}

	setTimeZone()

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
		Date:         data.Date.In(argTimeZone).Add(-time.Hour * 3),
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

	go d.geofenceService.ProcessData(data)

	return nil
}
