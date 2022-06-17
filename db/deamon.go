package db

import (
	"github.com/Fonzeca/Trackin/db/model"
	"github.com/Fonzeca/Trackin/db/query"
	"github.com/Fonzeca/Trackin/rest/json"
)

func Deamon(canal chan json.SimplyData) {
	for {
		data := <-canal
		processData(data)
	}
}

func processData(data json.SimplyData) {
	db, err := ObtenerConexionDb()
	if err != nil {
		//TODO: log error
		return
	}
	log := model.Log{
		Imei:         data.Imei,
		ProtocolType: data.ProtocolType,
		Latitud:      data.Latitude,
		Longitud:     data.Longitude,
		Date:         data.Date.Local(),
		Speed:        data.Speed,
		DeviceTemp:   data.DeviceTemp,
		Mileage:      data.Mileage,
		IsGps:        data.GpsWorking,
		IsHistory:    data.IsHistory,
		EngineStatus: data.EngineStatus,
	}

	q := query.Use(db).Log

	err = q.Create(&log)
	if err != nil {
		return
	}

}
