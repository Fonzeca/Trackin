package manager

import (
	"github.com/Fonzeca/Trackin/db"
	"github.com/Fonzeca/Trackin/db/model"
	"github.com/Fonzeca/Trackin/services"
)

type Manager struct {
	//Variable para setear ids para las vistas
	id int32
}

func NewManager() Manager {
	return Manager{id: 0}
}

func (ma *Manager) getId() int32 {
	ma.id++
	return ma.id
}

func (ma *Manager) GetLastLogByImei(imei string) (model.LastLogView, error) {
	// Usamos la nueva función con lock para evitar consultas duplicadas
	log, wasInCache := services.GetCachedPointsWithLock(imei, func() *model.Log {
		logResult := &model.Log{}
		db.DB.Select("imei", "latitud", "longitud", "speed", "date").Order("date desc").Where("imei = ?", imei).First(logResult)

		if logResult.Imei == "" {
			return nil
		}
		return logResult
	})

	// Solo actualizamos el caché en segundo plano si no estaba previamente cacheado
	if !wasInCache && log != nil {
		services.SetCachedPoints(imei, log)
	}

	var lastLog model.LastLogView
	if log != nil {
		lastLog = model.LastLogView{
			Imei: log.Imei,
			Location: model.Location{
				Latitutd: log.Latitud,
				Longitud: log.Longitud,
			},
			Speed: log.Speed,
			Date:  log.Date,
		}
	} else {
		lastLog = model.LastLogView{
			Imei: imei,
		}
	}

	return lastLog, nil
}

func (ma *Manager) GetVehiclesStateByImeis(only string, imeis model.ImeisBody) ([]model.StateLogView, error) {
	logs := []model.Log{}
	for _, imei := range imeis.Imeis {
		// Usamos la nueva función con lock para evitar consultas duplicadas
		log, wasInCache := services.GetCachedPointsWithLock(imei, func() *model.Log {
			return queryLogFromDB(imei, only)
		})

		if log != nil {
			logs = append(logs, *log)
		} else {
			// Si no se encontró nada, agregamos un log vacío
			logs = append(logs, model.Log{Imei: imei})
		}

		// Solo actualizamos el caché en segundo plano si no estaba previamente cacheado
		if !wasInCache && log != nil {
			go services.SetCachedPoints(imei, log)
		}
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
			Date:         log.Date,
		})
	}

	return stateLogsView, nil
}

func (ma *Manager) GetRouteByImei(requestRoute model.RouteRequest) ([]interface{}, error) {
	logs := []model.Log{}
	tx := db.DB.Select("date", "latitud", "longitud", "speed", "mileage", "engine_status", "azimuth").Where("imei = ? AND date BETWEEN ? AND ?", requestRoute.Imei, requestRoute.From, requestRoute.To).Order("date ASC").Find(&logs)

	routes := []interface{}{}
	movingData := []model.RouteDataView{}

	fromHour := ""
	fromDate := ""

	var isInStop bool = false
	var isMoving bool = false

	var initialMileage int32 = 0

	for index, log := range logs {
		if !log.EngineStatus {

			if isMoving {
				isMoving = false
				saveMovingLog(index-1, fromDate, fromHour, ma.getId(), &routes, &logs, movingData, initialMileage)
				movingData = []model.RouteDataView{}
			}

			if !isInStop {
				isInStop = true
				fromHour = log.Date.Format("15:04")
				fromDate = log.Date.Format("2006-01-02")
			}

			if index >= len(logs)-1 {
				saveStopLog(index, fromDate, fromHour, ma.getId(), &routes, &logs)
			}
			continue
		}

		if isInStop {
			isInStop = false
			saveStopLog(index-1, fromDate, fromHour, ma.getId(), &routes, &logs)
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
			Speed:     log.Speed,
			Azimuth:   log.Azimuth,
			Timestamp: log.Date.UnixMilli(),
		})

		if index >= len(logs)-1 {
			saveMovingLog(index, fromDate, fromHour, ma.getId(), &routes, &logs, movingData, initialMileage)
			movingData = []model.RouteDataView{}
		}
	}
	ma.id = 0
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

// queryLogFromDB ejecuta la consulta a la base de datos para obtener el último log
func queryLogFromDB(imei, only string) *model.Log {
	log := &model.Log{}
	if only != "" {
		db.DB.Select("imei", only, "date").Where("imei = ?", imei).Order("date DESC").First(log)
	} else {
		db.DB.Select("imei", "latitud", "longitud", "engine_status", "azimuth", "date").Where("imei = ?", imei).Order("date DESC").First(log)
	}

	// Si no se encontró nada, retornamos nil
	if log.Imei == "" {
		return nil
	}

	return log
}
