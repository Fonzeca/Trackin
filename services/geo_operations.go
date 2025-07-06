package services

import (
	"github.com/Fonzeca/Trackin/db/model"
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
