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

	if timeDifference == 0 || distance == 0 {
		return false // Avoid division by zero if the time difference or distance is zero
	}

	speed := distance / timeDifference

	return speed <= thresholdKmh
}

// ParseZoneToLoop converts a single zone into an S2 loop
func ParseZoneToLoop(zone model.ZoneView) (*s2.Loop, error) {
	if zone.Puntos == "" {
		return nil, fmt.Errorf("zone has no points")
	}

	pointsStr := strings.Split(zone.Puntos, ";")
	points := make([]s2.Point, 0, len(pointsStr))

	for _, pointStr := range pointsStr {
		point, err := getPointFromString(pointStr)
		if err != nil {
			// Log the error and continue to the next point
			fmt.Printf("Error parsing point '%s': %v\n", pointStr, err)
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
func ParseZonesToLoops(zones []model.ZoneView) ([]*s2.Loop, error) {
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
func ParseZonesToPolygon(zones []model.ZoneView) (*s2.Polygon, error) {
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
