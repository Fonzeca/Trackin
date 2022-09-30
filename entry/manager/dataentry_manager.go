package manager

import (
	"github.com/Fonzeca/Trackin/db"
	"github.com/Fonzeca/Trackin/db/model"
	"github.com/Fonzeca/Trackin/db/query"
	"github.com/Fonzeca/Trackin/entry/json"
)

type DataEntryManager struct {
	CanalEntrada chan json.SimplyData
}

func NewDataEntryManager() *DataEntryManager {
	instance := &DataEntryManager{
		CanalEntrada: make(chan json.SimplyData),
	}

	go instance.run()

	return instance
}

//Goroutine daemon
func (d *DataEntryManager) run() {
	for {
		data := <-d.CanalEntrada
		processData(data)
	}
}

func processData(data json.SimplyData) {
	//Evitamos datos inecesarios que llegan por equivocacion.
	if data.Latitude == 0 {
		return
	}

	db, close, err := db.ObtenerConexionDb()
	if err != nil {
		//TODO: log error
		return
	}
	defer close()

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

	err = q.Create(&log)
	if err != nil {
		return
	}

}
