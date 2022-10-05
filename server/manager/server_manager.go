package manager

import (
	"strconv"

	"github.com/Fonzeca/Trackin/db"
	"github.com/Fonzeca/Trackin/db/model"
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

func (ma *Manager) GetVehiclesStateByImeis(only string, imeis model.ImeisBody) ([]model.StateLogView, error) {
	db, close, err := db.ObtenerConexionDb()
	defer close()

	if err != nil {
		return nil, err
	}

	logs := []model.Log{}
	for _, imei := range imeis.Imeis {
		log := model.Log{}
		if only != "" {
			db.Select("imei", only, "date").Where("imei = ?", imei).Order("date DESC").Find(&log)
		} else {
			db.Select("imei", "latitud", "longitud", "engine_status", "azimuth", "date").Where("imei = ?", imei).Order("date DESC").Find(&log)
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

	return stateLogsView, nil
}

func (ma *Manager) GetRouteByImei(requestRoute model.RouteRequest) ([]interface{}, error) {
	db, close, err := db.ObtenerConexionDb()
	defer close()

	if err != nil {
		return nil, err
	}

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
			Speed:   log.Speed,
			Azimuth: log.Azimuth,
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

func (ma *Manager) GetZonesByEmpresaId(idParam string) ([]model.ZoneRequest, error) {
	db, close, err := db.ObtenerConexionDb()
	defer close()

	if err != nil {
		return nil, err
	}

	id, idParseErr := strconv.Atoi(idParam)

	if idParseErr != nil {
		return []model.ZoneRequest{}, idParseErr
	}

	zones := []model.Zona{}

	db.Find(&zones).Where("empresa_id = ?", id)
	tx := db.Model(&model.Zona{}).Joins("join zona_vehiculos on zona.id = zona_vehiculos.zona_id").Scan(&model.ZoneRequest{})

	zonesRequests := []model.ZoneRequest{}
	for _, zone := range zonesRequests {
		zonesRequests = append(zonesRequests, model.ZoneRequest{
			EmpresaId:     zone.EmpresaId,
			ColorLinea:    zone.ColorLinea,
			ColorRelleno:  zone.ColorRelleno,
			Puntos:        zone.Puntos,
			Nombre:        zone.Nombre,
			VehiculosIds:  zone.VehiculosIds,
			AvisarEntrada: zone.AvisarEntrada,
			AvisarSalida:  zone.AvisarSalida,
		})
	}
	return zonesRequests, tx.Error
}

func (ma *Manager) CreateZone(zoneRequest model.ZoneRequest) error {
	db, close, err := db.ObtenerConexionDb()
	defer close()

	if err != nil {
		return err
	}

	zone := model.Zona{
		EmpresaID:    int32(zoneRequest.EmpresaId),
		ColorLinea:   zoneRequest.ColorLinea,
		ColorRelleno: zoneRequest.ColorRelleno,
		Puntos:       zoneRequest.Puntos,
		Nombre:       zoneRequest.Nombre,
	}

	tx := db.Create(&zone)

	if tx.Error != nil {
		return tx.Error
	}

	zoneVehicles := []model.ZonaVehiculo{}
	for _, vehicleId := range zoneRequest.VehiculosIds {
		zoneVehicles = append(zoneVehicles, model.ZonaVehiculo{
			ZonaID:        zone.ID,
			VehiculoID:    int32(vehicleId),
			AvisarEntrada: zoneRequest.AvisarEntrada,
			AvisarSalida:  zoneRequest.AvisarSalida,
		})
	}

	tx = db.Create(&zoneVehicles)

	return tx.Error
}

func (ma *Manager) EditZoneById(idParam string, zoneView model.ZoneRequest) error {
	db, close, err := db.ObtenerConexionDb()
	defer close()

	if err != nil {
		return err
	}

	id, idParseErr := strconv.Atoi(idParam)

	if idParseErr != nil {
		return idParseErr
	}

	zone := model.Zona{
		ID:           int32(id),
		EmpresaID:    int32(zoneView.EmpresaId),
		ColorLinea:   zoneView.ColorLinea,
		ColorRelleno: zoneView.ColorRelleno,
		Puntos:       zoneView.Puntos,
		Nombre:       zoneView.Nombre,
	}
	tx := db.Save(&zone)

	return tx.Error
}

func (ma *Manager) DeleteZoneById(idParam string) error {
	db, close, err := db.ObtenerConexionDb()
	defer close()

	if err != nil {
		return err
	}

	id, idParseErr := strconv.Atoi(idParam)

	if idParseErr != nil {
		return idParseErr
	}

	zone := model.Zona{ID: int32(id)}
	tx := db.Delete(&zone)

	return tx.Error
}
