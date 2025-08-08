package manager

import (
	"errors"
	"log"
	"math"
	"time"

	"github.com/Fonzeca/Trackin/internal/core/services"
	db "github.com/Fonzeca/Trackin/internal/infrastructure/database"
	"github.com/Fonzeca/Trackin/internal/infrastructure/database/model"
	"github.com/Fonzeca/Trackin/internal/infrastructure/geolocation"
	"github.com/golang/geo/s2"
)

type routesManager struct {
	//Variable para setear ids para las vistas
	id           int32
	zonasManager IZonasManager
}

func newRoutesManager() IRoutesManager {
	return &routesManager{id: 0}
}

// SetZonasManager inyecta la dependencia del zonas manager
func (ma *routesManager) SetZonasManager(zonasManager IZonasManager) {
	ma.zonasManager = zonasManager
}

func (ma *routesManager) getId() int32 {
	ma.id++
	return ma.id
}

func (ma *routesManager) GetLastLogByImei(imei string) (model.LastLogView, error) {
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

func (ma *routesManager) GetVehiclesStateByImeis(imeis model.ImeisBody) ([]model.StateLogView, error) {
	logs := []model.Log{}
	for _, imei := range imeis.Imeis {
		// Usamos la nueva función con lock para evitar consultas duplicadas
		log, wasInCache := services.GetCachedPointsWithLock(imei, func() *model.Log {
			return queryLogFromDB(imei)
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

func (ma *routesManager) GetRouteByImeiAndZones(requestRoute model.RouteRequest, zones []model.ZoneView) ([]model.GpsRouteData, error) {
	// Esta función no está implementada en el código original
	// Aquí podrías implementar la lógica para obtener rutas por IMEI y zonas
	log.Println("GetRouteByImeiAndZones is not implemented")

	return nil, nil
}

func (ma *routesManager) GetRouteByImei(requestRoute model.RouteRequest) ([]model.GpsRouteData, error) {
	logs := []model.Log{}
	tx := db.DB.Select("date", "latitud", "longitud", "speed", "mileage", "engine_status", "azimuth").
		Where("imei = ? AND date BETWEEN ? AND ?", requestRoute.Imei, requestRoute.From, requestRoute.To).
		Order("date ASC").
		Find(&logs)

	routes := []model.GpsRouteData{}
	movingData := []model.GpsPoint{}

	fromHour := ""
	fromDate := ""

	var isInStop bool = false
	var isMoving bool = false

	var initialMileage int32 = 0

	// Usar el método público de la interfaz
	// logs = ma.CleanUpRouteBySpeedAnomaly(logs) // Comentado por ahora hasta ajustar tipos

	for index, log := range logs {
		if !log.EngineStatus {

			if isMoving {
				isMoving = false
				saveMovingLog(index-1, fromDate, fromHour, ma.getId(), &routes, &logs, movingData, initialMileage)
				movingData = []model.GpsPoint{}
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

			if index+1 < len(logs) {
				if !logs[index+1].EngineStatus &&
					logs[index+1].Latitud == log.Latitud &&
					logs[index+1].Longitud == log.Longitud {
					// Si el siguiente log también es una parada, no guardamos el log de parada
					continue
				}
			}

			isInStop = false
			saveStopLog(index-1, fromDate, fromHour, ma.getId(), &routes, &logs)
		}

		if !isMoving {
			isMoving = true
			fromHour = log.Date.Format("15:04")
			fromDate = log.Date.Format("2006-01-02")
			initialMileage = log.Mileage
		}

		movingData = append(movingData, model.GpsPoint{
			Azimuth:   log.Azimuth,
			Latitud:   log.Latitud,
			Longitud:  log.Longitud,
			Speed:     log.Speed,
			Timestamp: log.Date.UnixMilli(),
		})

		if index >= len(logs)-1 {
			saveMovingLog(index, fromDate, fromHour, ma.getId(), &routes, &logs, movingData, initialMileage)
			movingData = []model.GpsPoint{}
		}
	}
	ma.id = 0
	return routes, tx.Error
}

func (ma *routesManager) GetSummaryRoutesAndZones(request *model.SummaryRequest) ([]*model.PointIntersection, error) {

	if request.FromDate == 0 || request.ToDate == 0 || request.Imei == "" || request.EmpresaId == 0 {
		return nil, errors.New("not valid request")
	}

	logs, err := ma.GetLogsFromImeiAndDate(request)
	if err != nil {
		log.Println("Error getting logs:", err)
		return nil, err
	}

	if len(logs) == 0 {
		log.Println("No logs found for IMEI:", request.Imei)
		return nil, errors.New("no logs found")
	}

	zones, err := ma.zonasManager.GetZonesByEmpresaId(request.EmpresaId)
	if err != nil {
		log.Println("Error getting zones:", err)
		return nil, err
	}

	log.Printf("Intersecting %d logs with %d zones", len(logs), len(zones))

	// Map zone ID to S2 loop
	loopMap := make(map[int32]*s2.Loop)
	for _, zone := range zones {
		loop, err := geolocation.ParseZoneToLoop(zone)
		if err != nil {
			log.Printf("Error creating loop for zone %d: %v", zone.ID, err)
			continue
		}
		loopMap[zone.ID] = loop
	}

	zoneMap := make(map[int32]*model.Zona)
	for _, zone := range zones {
		z := zone // Make a copy to avoid referencing the loop variable
		zoneMap[zone.ID] = &z
	}

	// Intersect logs with zones
	intersections, err := geolocation.IntersectLogsAndZones(logs, zoneMap, loopMap)
	if err != nil {
		log.Printf("Error intersecting logs with zones: %v", err)
	}

	return intersections, nil
}

func (ma *routesManager) GetLogsFromImeiAndDate(request *model.SummaryRequest) ([]model.Log, error) {
	logs := []model.Log{}

	if request.FromDate == 0 || request.ToDate == 0 || request.Imei == "" || request.EmpresaId == 0 {
		return nil, errors.New("not valid request")
	}
	FromDateTime := time.UnixMilli(request.FromDate)
	ToDateTime := time.UnixMilli(request.ToDate)

	tx := db.DB.Select("imei", "date", "latitud", "longitud", "speed", "mileage", "engine_status", "azimuth").
		Where("imei = ? AND date BETWEEN ? AND ?", request.Imei, FromDateTime, ToDateTime).
		Order("date ASC").
		Find(&logs)

	if tx.Error != nil {
		log.Println("Error getting logs:", tx.Error)
		return nil, tx.Error
	}

	if len(logs) == 0 {
		log.Println("No logs found for IMEI:", request.Imei)
		return nil, errors.New("no logs found")
	}

	//filter logs when the car is stopped but keep only the first log of the stop
	filtered := logs[:1] // Always keep the first log
	for i := 1; i < len(logs); i++ {
		currentLog := logs[i]
		previousLog := logs[i-1]
		if !(!currentLog.EngineStatus && !previousLog.EngineStatus && currentLog.Latitud == previousLog.Latitud && currentLog.Longitud == previousLog.Longitud) {
			filtered = append(filtered, logs[i])
		}
	}
	logs = filtered

	return logs, nil
}

func saveStopLog(index int, fromDate string, fromHour string, id int32, routes *[]model.GpsRouteData, logs *[]model.Log) {
	*routes = append(*routes, model.GpsRouteData{
		Id:       id,
		Type:     "Parada",
		FromDate: fromDate,
		ToDate:   (*logs)[index].Date.Format("2006-01-02"),
		FromHour: fromHour,
		ToHour:   (*logs)[index].Date.Format("15:04"),
		Latitud:  (*logs)[index].Latitud,
		Longitud: (*logs)[index].Longitud,
		Data:     nil,
		Km:       0,
	})
}

func saveMovingLog(index int, fromDate string, fromHour string, id int32, routes *[]model.GpsRouteData, logs *[]model.Log, movingData []model.GpsPoint, initialMileage int32) {
	*routes = append(*routes, model.GpsRouteData{
		Id:       id,
		Type:     "Viaje",
		FromDate: fromDate,
		ToDate:   (*logs)[index].Date.Format("2006-01-02"),
		FromHour: fromHour,
		ToHour:   (*logs)[index].Date.Format("15:04"),
		Km:       ((*logs)[index].Mileage - initialMileage) / 1000,
		Data:     movingData,
	})
}

// queryLogFromDB ejecuta la consulta a la base de datos para obtener el último log
func queryLogFromDB(imei string) *model.Log {
	log := &model.Log{}
	db.DB.Raw("SELECT imei, latitud, longitud, engine_status, azimuth, date FROM log WHERE imei = ? ORDER BY date DESC LIMIT 1", imei).
		Scan(log)

	// Si no se encontró nada, retornamos nil
	if log.Imei == "" {
		return nil
	}

	return log
}

// distanceOf2Points calcula la distancia entre dos puntos GPS usando la fórmula de Haversine
// Retorna la distancia en metros
func distanceOf2Points(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusMeters = 6371000 // Radio de la Tierra en metros

	// Convertir grados a radianes
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLatRad := (lat2 - lat1) * math.Pi / 180
	deltaLonRad := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLatRad/2)*math.Sin(deltaLatRad/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLonRad/2)*math.Sin(deltaLonRad/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusMeters * c
}

// cleanUpRouteBySpeedAnomaly limpia una ruta eliminando puntos con anomalías de velocidad
// CleanUpRouteBySpeedAnomaly limpia una ruta eliminando puntos con anomalías de velocidad
func (ma *routesManager) CleanUpRouteBySpeedAnomaly(route []model.GpsPoint) []model.GpsPoint {
	cleanedRoute := []model.GpsPoint{}
	const speedThreshold = 350.0 // Umbral de velocidad en km/h

	if len(route) <= 1 {
		return route
	}

	// Agregar el primer punto siempre
	cleanedRoute = append(cleanedRoute, route[0])

	for i := 1; i < len(route); i++ {
		currentPoint := route[i]
		previousPoint := route[i-1]

		// Saltar si los timestamps son iguales
		if currentPoint.Timestamp == previousPoint.Timestamp {
			continue
		}

		// Calcular distancia entre puntos en metros
		distanceOfPointsMeters := distanceOf2Points(
			previousPoint.Latitud, previousPoint.Longitud,
			currentPoint.Latitud, currentPoint.Longitud,
		)

		// Convertir a kilómetros
		distanceOfPoints := distanceOfPointsMeters / 1000.0

		// Calcular diferencia de tiempo
		timeDifferenceMilliseconds := currentPoint.Timestamp - previousPoint.Timestamp
		timeDiffHours := float64(timeDifferenceMilliseconds) / 3600000.0 // Convertir a horas

		// Calcular velocidad en km/h
		speed := distanceOfPoints / timeDiffHours

		// Verificar si hay anomalía de velocidad
		if speed > speedThreshold {
			log.Printf("Anomaly detected: Speed %.2f km/h between points with timestamps %d and %d",
				speed,
				currentPoint.Timestamp,
				previousPoint.Timestamp)
			continue
		}

		cleanedRoute = append(cleanedRoute, currentPoint)
	}

	return cleanedRoute
}
