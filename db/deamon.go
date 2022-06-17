package db

import (
	"fmt"
	"time"

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
	fmt.Println(time.Local.String())

	log := model.Log{
		Imei:     data.Imei,
		Latitud:  data.Latitude,
		Longitud: data.Longitude,
		Date:     data.Date.Local(),
	}

	q := query.Use(db).Log

	err = q.Create(&log)
	if err != nil {
		return
	}

}
