package manager

import (
	"time"

	"github.com/Fonzeca/Trackin/db"
	"github.com/Fonzeca/Trackin/db/model"
	"github.com/Fonzeca/Trackin/db/query"
	"github.com/Fonzeca/Trackin/entry/json"
)

var argTimeZone *time.Location

type DataEntryManager struct {
	CanalEntrada chan json.SimplyData
}

func setTimeZone() {
	arg, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		panic(err)
	}
	argTimeZone = arg
}

func NewDataEntryManager() *DataEntryManager {
	instance := &DataEntryManager{
		CanalEntrada: make(chan json.SimplyData),
	}

	go instance.run()

	setTimeZone()

	return instance
}

// Goroutine daemon
func (d *DataEntryManager) run() {
	for {
		data := <-d.CanalEntrada
		d.ProcessData(data)
	}
}

func (d *DataEntryManager) ProcessData(data json.SimplyData) error {
	//Evitamos datos inecesarios que llegan por equivocacion.
	if data.Latitude == 0 {
		return nil
	}

	db, close, err := db.ObtenerConexionDb()
	if err != nil {
		//TODO: log error
		return err
	}
	defer close()

	log := model.Log{
		Imei:         data.Imei,
		ProtocolType: data.ProtocolType,
		Latitud:      data.Latitude,
		Longitud:     data.Longitude,
		Date:         data.Date.In(argTimeZone),
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

	err = q.Create(&log)
	if err != nil {
		return err
	}
	return nil
}
