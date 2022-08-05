package manager

import (
	"github.com/Fonzeca/Trackin/db"
	"github.com/Fonzeca/Trackin/db/model"
)

type Manager struct {
}

func NewManager() Manager {
	return Manager{}
}

func (ma *Manager) GetLastLogByImei(imei string) (model.LastLogView, error) {
	db, err := db.ObtenerConexionDb()

	if err != nil {
		return model.LastLogView{}, err
	}

	log := model.Log{}
	tx := db.Select("imei", "latitud", "longitud", "speed", "date").Last(&log, "imei = ?", imei)

	lastLog := model.LastLogView{
		Imei:     log.Imei,
		Latitutd: log.Latitud,
		Longitud: log.Longitud,
		Speed:    log.Speed,
		Date:     log.Date,
	}

	return lastLog, tx.Error
}
