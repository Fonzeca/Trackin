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

func ParseZonesToPolygon(zones []model.ZoneView) (*s2.Polygon, error) {
	if len(zones) == 0 {
		return nil, fmt.Errorf("no zones provided")
	}

	var allLoops []*s2.Loop

	// Iterate through each zone and convert the points to a polygon

	for _, zone := range zones {
		if zone.Puntos != "" {
			pointsStr := strings.Split(zone.Puntos, ";")
			points := make([]s2.Point, 0, len(pointsStr))

			for _, pointStr := range pointsStr {
				point, err := getPointFromString(pointStr)
				if err != nil {
					continue
				}
				points = append(points, *point)
			}

			// Ensure the points are in the correct order and form a closed loop
			// Only create a polygon if there are enough points
			if len(points) > 1 {
				loop := s2.LoopFromPoints(points)

				loop.Normalize() // Normalize the loop to ensure it is valid

				allLoops = append(allLoops, loop)
			}
		}
	}

	if len(allLoops) == 0 {
		return nil, fmt.Errorf("no valid zones found")
	}

	return s2.PolygonFromLoops(allLoops), nil
}

func getPointFromString(pointStr string) (*s2.Point, error) {
	coords := strings.Split(pointStr, ",")
	if len(coords) != 2 {
		return nil, fmt.Errorf("invalid point format: %s", pointStr)
	}

	lat, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return nil, err
	}

	lng, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return nil, err
	}

	latLng := s2.LatLngFromDegrees(lat, lng)
	point := s2.PointFromLatLng(latLng)

	return &point, nil
}
