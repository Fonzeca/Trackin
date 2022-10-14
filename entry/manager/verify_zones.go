package manager

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	jsonEncoder "encoding/json"

	"github.com/Fonzeca/Trackin/db/model"
	"github.com/Fonzeca/Trackin/entry/json"
	"github.com/Fonzeca/Trackin/server/manager"
	"github.com/Fonzeca/Trackin/services"
	"github.com/rabbitmq/amqp091-go"
)

type Point struct {
	lat float64
	lng float64
}

type GeofenceDetector struct {
	last_logs map[string]Point
	manager   manager.ZonasManager
}

func NewGeofenceDetector() *GeofenceDetector {
	zm := manager.NewZonasManager()
	return &GeofenceDetector{
		manager:   *zm,
		last_logs: make(map[string]Point),
	}
}

func (d *GeofenceDetector) ProcessData(data json.SimplyData) error {

	fmt.Println("ZonesManager: Procesando " + data.Imei)

	imei := data.Imei

	//Punto del imei viejo
	oldVehiclePoint, ok := d.last_logs[imei]

	//Punto del imei nuevo
	var currentVehiclePoint Point = Point{lat: data.Latitude, lng: data.Longitude}

	//Verificamos si existe el punto viejo
	if !ok {
		d.last_logs[imei] = currentVehiclePoint
		return nil
	}

	//Obtenemos la config del imei con la zona
	zonesConfig, err := d.manager.GetZoneConfigByImei(imei)
	if err != nil {
		return err
	}

	if len(zonesConfig) <= 0 {
		//El vehiculo no está asociado a ninguna zona, por lo tanto no hacemos nada
		return nil
	}

	//Recorremos cada zona que tiene asignado el imei
	for _, zoneConfig := range zonesConfig {

		//Si el vehíuclo no tiene que avisar nada en una zona, salteamos la zona
		if !zoneConfig.AvisarEntrada || !zoneConfig.AvisarSalida {
			continue
		}

		//Obtenemos los puntos del poligono de la zona, porque estan almacenados como un string
		polygon, err := getPolygonFromString(zoneConfig.Puntos)
		if err != nil {
			return err
		}

		//Verificamos los 2 puntos, si estan adentro de la zona
		var isCurrentVehiclePointInZone bool = isPointInPolygon(currentVehiclePoint, polygon)
		var isOldVehiclePointInZone bool = isPointInPolygon(oldVehiclePoint, polygon)

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
		} else if zoneConfig.AvisarSalida {
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
			services.GlobalChannel.PublishWithContext(context.Background(), "carmind", "notification.zone.back.preparing", true, true, amqp091.Publishing{
				ContentType: "application/json",
				Body:        zoneNotificationBytes,
			})
			fmt.Println("Mensaje mandando:" + imei)
		}

	}

	return nil
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

	var isInside bool = false

	j := len(polygon) - 1
	for i := 0; i < len(polygon); j = i + 1 {
		if (polygon[i].lng > p.lng) != (polygon[j].lng > p.lng) &&
			p.lat < (polygon[j].lat-polygon[i].lat)*(p.lng-polygon[i].lng)/(polygon[j].lng-polygon[i].lng)+polygon[i].lat {
			isInside = !isInside
		}
	}

	return isInside
}
