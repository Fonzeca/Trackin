package core

import (
	"fmt"
	"strconv"

	"github.com/Fonzeca/Trackin/internal/infrastructure/database/model"
)

// MapIntersectionsToCSV maps the intersections to a CSV format
func MapIntersectionsToCSV(intersections []*model.PointIntersection) ([][]string, error) {

	data := make([][]string, 0)
	data = append(data, []string{
		"IMEI",
		"Fecha",
		"Latitud",
		"Longitud",
		"Velocidad (km/h)",
		"Zonas",
		"Velocidad Maxima",
		"Exceso de Velocidad",
		"Motor encendido",
		"Link maps",
	})

	// Write data
	for _, intersection := range intersections {
		log := intersection.Log

		// Get zone names using helper function
		zonesStr := getZoneNames(intersection.Zones)

		// Get the lowest speed limit using helper function
		minSpeedLimit := getMinSpeedLimit(intersection.Zones)

		// Check speed excess
		speedExcess := "No"
		if minSpeedLimit > 0 && float64(log.Speed) > minSpeedLimit {
			speedExcess = "Sí"
		}

		minSpeedStr := ""
		if minSpeedLimit > 0 {
			minSpeedStr = strconv.FormatFloat(minSpeedLimit, 'f', 2, 64)
		}

		record := []string{
			log.Imei,
			log.Date.Format("02/01/2006 15:04:05"),
			"lat " + strconv.FormatFloat(log.Latitud, 'f', 6, 64),
			"lng " + strconv.FormatFloat(log.Longitud, 'f', 6, 64),
			strconv.FormatFloat(float64(log.Speed), 'f', 2, 32),
			zonesStr,
			minSpeedStr,
			speedExcess,
			fmt.Sprintf("%t", log.EngineStatus),
			fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%f,%f", log.Latitud, log.Longitud),
		}

		data = append(data, record)
	}

	return data, nil
}

func getZoneNames(zones []*model.Zona) string {
	if len(zones) == 0 {
		return ""
	}

	var zoneNames []string
	for _, zone := range zones {
		zoneNames = append(zoneNames, zone.Nombre)
	}

	zonesStr := zoneNames[0]
	for i := 1; i < len(zoneNames); i++ {
		zonesStr += "; " + zoneNames[i]
	}

	return zonesStr
}

// getMinSpeedLimit finds the lowest speed limit among the zones
func getMinSpeedLimit(zones []*model.Zona) float64 {
	minSpeedLimit := -1.0 // -1 indicates no limit

	for _, zone := range zones {
		if zone.VelocidadMaxima > 0 {
			if minSpeedLimit == -1 || zone.VelocidadMaxima < minSpeedLimit {
				minSpeedLimit = zone.VelocidadMaxima
			}
		}
	}

	return minSpeedLimit
}
