package geolocation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Fonzeca/Trackin/internal/infrastructure/database/model"
	"github.com/golang/geo/s2"
)

const thresholdKmh = 350.0    // km/h threshold for speed filtering
const earthRadiusKm = 6371.01 // Average radius of the Earth in km

// IsValidPoint checks if the transition from current to next log entry is valid based on speed and time constraints.
// It ensures that the current log's date is before the next log's date, calculates the distance between the two points,
// and checks if the speed is within the defined threshold.
func IsValidPoint(current, next *model.Log) bool {
	if current == nil || next == nil {
		return false
	}

	if current.Date.After(next.Date) {
		return false // Ensure the current date is before the next date
	}

	currentPoint := s2.LatLngFromDegrees(current.Latitud, current.Longitud)
	nextPoint := s2.LatLngFromDegrees(next.Latitud, next.Longitud)

	distance := nextPoint.Distance(currentPoint).Radians() * earthRadiusKm // Calculate the distance between the two points in km

	// Calculate the time difference in hours
	timeDifference := next.Date.Sub(current.Date).Hours()

	if timeDifference == 0 {
		return false // Avoid division by zero if the time difference is zero
	}

	const minDistanceThresholdKm = 0.005 // 5 metros mínimo
	const minTimeThresholdMinutes = 1    // 1 minuto mínimo

	if distance <= minDistanceThresholdKm {
		if timeDifference*60 >= minTimeThresholdMinutes {
			return true
		}
		return false
	}

	speed := distance / timeDifference

	return speed <= thresholdKmh
}

// IsValidPointFlexible checks if the transition from current to next log entry is valid based only on speed constraints.
// This is a more flexible version that only validates excessive speeds, useful for cleaning GPS data with intermittent issues.
func IsValidPointFlexible(current, next *model.Log) bool {
	if current == nil || next == nil {
		return false
	}

	currentPoint := s2.LatLngFromDegrees(current.Latitud, current.Longitud)
	nextPoint := s2.LatLngFromDegrees(next.Latitud, next.Longitud)

	searchingPoint := s2.LatLngFromDegrees(-49.303356, -67.744521)

	// want to verify if currentPoint are in radius of 100 meters from a new point
	if currentPoint.Distance(searchingPoint).Radians()*earthRadiusKm <= 100*0.001 {
		return false
	}

	// Allow points with same timestamp or reversed order (more flexible)
	if current.Date.Equal(next.Date) {
		return true // Same timestamp is acceptable
	}

	distance := nextPoint.Distance(currentPoint).Radians() * earthRadiusKm

	if distance > 20 {
		// Si la distancia es mayor a 20 km, consideramos que es inválido directamente
		return false
	}

	// Calculate absolute time difference in hours
	timeDifference := next.Date.Sub(current.Date)
	if timeDifference < 0 {
		timeDifference = -timeDifference // Make it absolute
	}
	timeDifferenceHours := timeDifference.Hours()

	// If time difference is zero, consider it valid (same timestamp)
	if timeDifferenceHours == 0 {
		return true
	}

	// More flexible speed threshold for debugging (e.g., 500 km/h instead of 350)
	const flexibleThresholdKmh = 350.0
	speed := distance / timeDifferenceHours

	return speed <= flexibleThresholdKmh
}

// ParseZoneToLoop converts a single zone into an S2 loop
func ParseZoneToLoop(zone model.Zona) (*s2.Loop, error) {
	if zone.Puntos == "" {
		return nil, fmt.Errorf("zone has no points")
	}

	pointsStr := strings.Split(zone.Puntos, ";")
	points := make([]s2.Point, 0, len(pointsStr))

	for _, pointStr := range pointsStr {
		point, err := getPointFromString(pointStr)
		if err != nil {
			// Skip invalid points and continue to the next point
			continue
		}
		points = append(points, *point)
	}

	// Only create a loop if there are enough points
	if len(points) <= 1 {
		return nil, fmt.Errorf("not enough valid points to create a loop")
	}

	loop := s2.LoopFromPoints(points)
	loop.Normalize() // Normalize the loop to ensure it is valid

	return loop, nil
}

// ParseZonesToLoops converts an array of zones into an array of S2 loops
func ParseZonesToLoops(zones []model.Zona) ([]*s2.Loop, error) {
	if len(zones) == 0 {
		return nil, fmt.Errorf("no zones provided")
	}

	var allLoops []*s2.Loop

	for _, zone := range zones {
		loop, err := ParseZoneToLoop(zone)
		if err != nil {
			// Continue with other zones if one fails
			continue
		}
		allLoops = append(allLoops, loop)
	}

	if len(allLoops) == 0 {
		return nil, fmt.Errorf("no valid zones found")
	}

	return allLoops, nil
}

// ParseZonesToPolygon converts an array of zones into a unified S2 polygon
func ParseZonesToPolygon(zones []model.Zona) (*s2.Polygon, error) {
	loops, err := ParseZonesToLoops(zones)
	if err != nil {
		return nil, err
	}

	return s2.PolygonFromLoops(loops), nil
}

func getPointFromString(pointStr string) (*s2.Point, error) {
	coords := strings.Split(pointStr, ",")
	if len(coords) != 2 {
		return nil, fmt.Errorf("invalid point format: %s", pointStr)
	}

	// Replace comma decimal separator with dot for both latitude and longitude
	latStr := strings.ReplaceAll(strings.TrimSpace(coords[0]), ",", ".")
	lngStr := strings.ReplaceAll(strings.TrimSpace(coords[1]), ",", ".")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return nil, err
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		return nil, err
	}

	latLng := s2.LatLngFromDegrees(lat, lng)
	point := s2.PointFromLatLng(latLng)

	return &point, nil
}

func IntersectLogsAndZones(logs []model.Log, zoneMap map[int32]*model.Zona, loopMap map[int32]*s2.Loop) ([]*model.PointIntersection, error) {
	if len(logs) == 0 || len(zoneMap) == 0 || len(loopMap) == 0 {
		return nil, fmt.Errorf("logs, zones or loops cannot be empty")
	}

	var intersections []*model.PointIntersection

	for _, logGeo := range logs {
		logGeoRef := logGeo // Create a copy to avoid referencing the loop variable
		intersec := &model.PointIntersection{
			Log:   &logGeoRef,
			Zones: make([]*model.Zona, 0),
		}

		for zoneId, loop := range loopMap {
			if loop.ContainsPoint(s2.PointFromLatLng(s2.LatLngFromDegrees(logGeo.Latitud, logGeo.Longitud))) {
				if zone, exists := zoneMap[zoneId]; exists {
					intersec.Zones = append(intersec.Zones, zone)
				}
			}
		}

		// Add intersection to array (even if it has no zones)
		intersections = append(intersections, intersec)
	}

	return intersections, nil
}
