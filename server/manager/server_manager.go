package manager

import (
	"time"

	"github.com/Fonzeca/Trackin/db"
	"github.com/Fonzeca/Trackin/db/model"
	"gorm.io/gorm"
)

type Manager struct {
}

func NewManager() Manager {
	return Manager{}
}

func (ma *Manager) GetLastLogByImei(imei string) (model.LastLogView, error) {
	db, close, err := db.ObtenerConexionDb()
	defer close()

	if err != nil {
		return model.LastLogView{}, err
	}

	log := model.Log{}
	tx := db.Select("imei", "latitud", "longitud", "speed", "date").Order("date desc").Where("imei = ?", imei).First(&log)

	lastLog := model.LastLogView{
		Imei: log.Imei,
		Location: model.Location{
			Latitutd: log.Latitud,
			Longitud: log.Longitud,
		},
		Speed: log.Speed,
		Date:  log.Date,
	}

	return lastLog, tx.Error
}

func (ma *Manager) GetVehiclesStateByImeis(only string, imeis model.Imeis) ([]model.StateLogView, error) {
	db, close, err := db.ObtenerConexionDb()
	defer close()

	if err != nil {
		return nil, err
	}

	logs := []model.Log{}
	var tx *gorm.DB
	if only != "" {
		tx = db.Select("imei", only, "max(date)").Where("imei IN ?", imeis.Imeis).Group("imei").Find(&logs)
	} else {
		tx = db.Select("imei", "latitud", "longitud", "engine_status", "azimuth", "max(date)").Where("imei IN ?", imeis.Imeis).Group("imei").Find(&logs)
	}

	stateLogsView := []model.StateLogView{}
	for _, log := range logs {
		stateLogsView = append(stateLogsView, model.StateLogView{
			Imei: log.Imei,
			Location: model.Location{
				Latitutd: log.Latitud,
				Longitud: log.Longitud,
			},
			EngineStatus: log.EngineStatus,
			Azimuth:      log.Azimuth,
		})
	}

	return stateLogsView, tx.Error
}

func (ma *Manager) GetRouteByImeiAndDate(imei string, from string, to string) (model.RouteView, error) {
	db, close, err := db.ObtenerConexionDb()
	defer close()

	if err != nil {
		return model.RouteView{}, err
	}

	var fromDate time.Time
	var toDate time.Time

	if from == "" || to == "" {
		fromDate = time.Now()
		fromDate = time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, fromDate.Location())

		toDate = time.Now().AddDate(0, 0, 1)
		toDate = time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 0, 0, 0, 0, toDate.Location())
		toDate = toDate.Add(-time.Second)
	}

	logs := []model.Log{}
	db.Table("log").Where("imei = ? AND date BETWEEN ? AND ?", imei, fromDate.In(time.UTC), toDate.In(time.UTC)).Order("date DESC").Find(&logs)

	view := model.RouteView{
		Imei: imei,
		From: fromDate.Format(time.RFC3339),
		To:   toDate.Format(time.RFC3339),
	}

	data := []model.RouteDataView{}
	for _, log := range logs {
		dataLog := model.RouteDataView{
			Date:  log.Date,
			Speed: log.Speed,
			Location: model.Location{
				Latitutd: log.Latitud,
				Longitud: log.Longitud,
			},
		}
		data = append(data, dataLog)
	}

	view.Data = data

	return view, nil
}
