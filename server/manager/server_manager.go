package manager

import (
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

func (ma *Manager) GetVehiclesStateByImeis(only string, imeis model.StateRequest) ([]model.StateLogView, error) {
	db, close, err := db.ObtenerConexionDb()
	defer close()

	if err != nil {
		return nil, err
	}

	logs := []model.Log{}
	var tx *gorm.DB
	for _, imei := range imeis.Imeis {
		log := model.Log{}
		if only != "" {
			tx = db.Select("imei", only, "date").Where("imei = ?", imei).Order("date DESC").Find(&log)
		} else {
			tx = db.Select("imei", "latitud", "longitud", "engine_status", "azimuth", "date").Where("imei = ?", imei).Order("date DESC").Find(&log)
		}
		logs = append(logs, log)
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

func (ma *Manager) GetRouteByImei(requestRoute model.RouteRequest) ([]interface{}, error) {
	db, close, err := db.ObtenerConexionDb()
	defer close()

	if err != nil {
		return nil, err
	}
	var id int32 = 0

	logs := []model.Log{}
	tx := db.Select("date", "latitud", "longitud", "speed", "mileage", "engine_status", "azimuth").Where("imei = ? AND date BETWEEN ? AND ?", requestRoute.Imei, requestRoute.From, requestRoute.To).Order("date ASC").Find(&logs)

	routes := []interface{}{}
	movingData := []model.RouteDataView{}

	fromHour := ""
	fromDate := ""

	var isInStop bool = false
	var isMoving bool = false

	var initialMileage int32 = 0

	for index, log := range logs {
		if !*log.EngineStatus {

			if isMoving {
				isMoving = false
				saveMovingLog(index-1, fromDate, fromHour, id, &routes, &logs, movingData, initialMileage)
				movingData = []model.RouteDataView{}
			}

			if !isInStop {
				isInStop = true
				fromHour = log.Date.Format("15:04")
				fromDate = log.Date.Format("2006-01-02")
			}

			if index >= len(logs)-1 {
				saveStopLog(index, fromDate, fromHour, id, &routes, &logs)
			}
			continue
		}

		if isInStop {
			isInStop = false
			saveStopLog(index-1, fromDate, fromHour, id, &routes, &logs)
		}

		if !isMoving {
			isMoving = true
			fromHour = log.Date.Format("15:04")
			fromDate = log.Date.Format("2006-01-02")
			initialMileage = log.Mileage
		}

		movingData = append(movingData, model.RouteDataView{
			Location: model.Location{
				Latitutd: log.Latitud,
				Longitud: log.Longitud,
			},
			Speed:   log.Speed,
			Azimuth: log.Azimuth,
		})

		if index >= len(logs)-1 {
			saveMovingLog(index, fromDate, fromHour, id, &routes, &logs, movingData, initialMileage)
			movingData = []model.RouteDataView{}
		}
		id++
	}
	return routes, tx.Error
}

func saveStopLog(index int, fromDate string, fromHour string, id int32, routes *[]interface{}, logs *[]model.Log) {
	*routes = append(*routes, model.StopView{
		RouteView: model.RouteView{
			Id:       id,
			Type:     "Parada",
			FromDate: fromDate,
			ToDate:   (*logs)[index].Date.Format("2006-01-02"),
			FromHour: fromHour,
			ToHour:   (*logs)[index].Date.Format("15:04"),
		},
		Location: model.Location{
			Latitutd: (*logs)[index].Latitud,
			Longitud: (*logs)[index].Longitud,
		},
	})
}

func saveMovingLog(index int, fromDate string, fromHour string, id int32, routes *[]interface{}, logs *[]model.Log, movingData []model.RouteDataView, initialMileage int32) {
	*routes = append(*routes, model.MoveView{
		RouteView: model.RouteView{
			Id:       id,
			Type:     "Viaje",
			FromDate: fromDate,
			ToDate:   (*logs)[index].Date.Format("2006-01-02"),
			FromHour: fromHour,
			ToHour:   (*logs)[index].Date.Format("15:04"),
		},
		KM:   ((*logs)[index].Mileage - initialMileage) / 1000,
		Data: movingData,
	})
}
