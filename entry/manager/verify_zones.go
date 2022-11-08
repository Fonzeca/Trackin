package manager

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	jsonEncoder "encoding/json"

	"github.com/Fonzeca/Trackin/db/model"
	"github.com/Fonzeca/Trackin/entry/json"
	"github.com/Fonzeca/Trackin/server/manager"
	"github.com/Fonzeca/Trackin/services"
)

type Point struct {
	lat  float64
	lng  float64
	date time.Time
}

type GeofenceDetector struct {
	last_logs map[string]Point

	imeisChannels map[string]chan json.SimplyData
	manager       manager.IZonasManager
}

func NewGeofenceDetector() *GeofenceDetector {
	zm := manager.ZonasManager
	return &GeofenceDetector{
		manager:       zm,
		last_logs:     make(map[string]Point),
		imeisChannels: make(map[string]chan json.SimplyData),
	}
}

func (d *GeofenceDetector) DispatchMessage(data json.SimplyData) {
	fmt.Printf("Data recived, imei %s\n", data.Imei)
	if data.EngineStatus {
		if _, ok := d.imeisChannels[data.Imei]; !ok {
			newChannel := make(chan json.SimplyData)
			d.imeisChannels[data.Imei] = newChannel
			go d.Worker(newChannel, data.Imei)
		}
		d.imeisChannels[data.Imei] <- data
	} else {
		if channel, ok := d.imeisChannels[data.Imei]; ok {
			data.CloseChannel = true
			channel <- data
		}
	}
}

func (d *GeofenceDetector) Worker(channel chan json.SimplyData, imei string) {
	defer func() {
		d.imeisChannels[imei] = nil
		delete(d.imeisChannels, imei)
		close(channel)
	}()

	//Obtenemos la config del imei con la zona
	zonesConfig, _ := d.manager.GetZoneConfigByImei(imei)
	go func() {
		for {
			time.Sleep(5 * time.Second)
			zonesConfig, _ = d.manager.GetZoneConfigByImei(imei)
		}
	}()

	var oldVehiclePoint *Point

	for {
		data := <-channel

		if zonesConfig != nil && len(zonesConfig) <= 0 {
			continue
		}

		var currentVehiclePoint *Point = &Point{lat: data.Latitude, lng: data.Longitude, date: data.Date}
		if oldVehiclePoint == nil {
			oldVehiclePoint = currentVehiclePoint
			continue
		}

		//Siempre los nuevos puntos tienen que ser de un date mayor al old
		if currentVehiclePoint.date.Before(oldVehiclePoint.date) {
			//Esto significa que es viejo
			continue
		}

		for _, zoneConfig := range zonesConfig {
			//Si el vehíuclo no tiene que avisar nada en una zona, salteamos la zona
			if !zoneConfig.AvisarEntrada || !zoneConfig.AvisarSalida {
				continue
			}

			//Obtenemos los puntos del poligono de la zona, porque estan almacenados como un string
			polygon, err := getPolygonFromString(zoneConfig.Puntos)
			if err != nil {
				fmt.Println(err)
				continue
			}

			//Verificamos los 2 puntos, si estan adentro de la zona
			var isCurrentVehiclePointInZone bool = isPointInPolygon(*currentVehiclePoint, polygon)
			var isOldVehiclePointInZone bool = isPointInPolygon(*oldVehiclePoint, polygon)

			var zoneNotification *model.ZoneNotification
			if zoneConfig.AvisarEntrada {
				if isCurrentVehiclePointInZone && !isOldVehiclePointInZone {
					fmt.Println("Entro! : " + imei)
					zoneNotification = &model.ZoneNotification{
						Imei:      imei,
						ZoneName:  zoneConfig.Nombre,
						ZoneID:    int(zoneConfig.Id),
						EventType: "entra",
					}
				}
			}

			if zoneConfig.AvisarSalida {
				if !isCurrentVehiclePointInZone && isOldVehiclePointInZone {
					fmt.Println("Salio! : " + imei)
					zoneNotification = &model.ZoneNotification{
						Imei:      imei,
						ZoneName:  zoneConfig.Nombre,
						ZoneID:    int(zoneConfig.Id),
						EventType: "sale",
					}
				}
			}

			if zoneNotification != nil {
				zoneNotificationBytes, _ := jsonEncoder.Marshal(zoneNotification)
				fmt.Println("Por manadar message:" + imei)
				err := services.GlobalSender.SendMessage(context.Background(), "notification.zone.back.preparing", zoneNotificationBytes)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Mensaje mandando:" + imei)
			}

		}

		oldVehiclePoint = currentVehiclePoint

		if data.CloseChannel {
			return
		}
	}
}

func getPolygonFromString(s string) ([]Point, error) {
	var splittedPoints []string = strings.Split(s, "; ")
	var polygon []Point

	for _, point := range splittedPoints {
		var splittedPoint []string = strings.Split(point, ",")
		lat, err := strconv.ParseFloat(splittedPoint[0], 32)
		if err != nil {
			return nil, err
		}
		lng, err := strconv.ParseFloat(splittedPoint[1], 32)
		if err != nil {
			return nil, err
		}
		polygon = append(polygon, Point{lat: lat, lng: lng})
	}
	return polygon, nil
}

func isPointInPolygon(p Point, polygon []Point) bool {

	minX := polygon[0].lat
	maxX := polygon[0].lat
	minY := polygon[0].lng
	maxY := polygon[0].lng

	for _, point := range polygon {
		minX = math.Min(point.lat, minX)
		maxX = math.Max(point.lat, maxX)
		minY = math.Min(point.lng, minY)
		maxY = math.Max(point.lng, maxY)
	}

	if p.lat < minX || p.lat > maxX || p.lng < minY || p.lng > maxY {
		return false
	}

	start := len(polygon) - 1
	end := 0

	contains := intersectsWithRaycast(&p, &polygon[start], &polygon[end])

	for i := 1; i < len(polygon); i++ {
		if intersectsWithRaycast(&p, &polygon[i-1], &polygon[i]) {
			contains = !contains
		}
	}

	return contains
}

// Using the raycast algorithm, this returns whether or not the passed in point
// Intersects with the edge drawn by the passed in start and end points.
// Original implementation: http://rosettacode.org/wiki/Ray-casting_algorithm#Go
func intersectsWithRaycast(point *Point, start *Point, end *Point) bool {
	// Always ensure that the the first point
	// has a y coordinate that is less than the second point
	if start.lng > end.lng {

		// Switch the points if otherwise.
		start, end = end, start

	}

	// Move the point's y coordinate
	// outside of the bounds of the testing region
	// so we can start drawing a ray
	for point.lng == start.lng || point.lng == end.lng {
		newLng := math.Nextafter(point.lng, math.Inf(1))
		point = &Point{lat: point.lat, lng: newLng}

	}

	// If we are outside of the polygon, indicate so.
	if point.lng < start.lng || point.lng > end.lng {
		return false
	}

	if start.lat > end.lat {
		if point.lat > start.lat {
			return false
		}
		if point.lat < end.lat {
			return true
		}

	} else {
		if point.lat > end.lat {
			return false
		}
		if point.lat < start.lat {
			return true
		}
	}

	raySlope := (point.lng - start.lng) / (point.lat - start.lat)
	diagSlope := (end.lng - start.lng) / (end.lat - start.lat)

	return raySlope >= diagSlope
}
