package manager

import (
	"github.com/Fonzeca/Trackin/entry/json"
	"github.com/Fonzeca/Trackin/server/manager"
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
	imei := data.Imei

	_, ok := d.last_logs[imei]

	if !ok {
		d.last_logs[imei] = Point{lat: data.Latitude, lng: data.Longitude}
		return nil
	}

	return nil
}
