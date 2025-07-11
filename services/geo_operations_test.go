package services

import (
	"testing"
	"time"

	"github.com/Fonzeca/Trackin/db/model"
	"github.com/golang/geo/s2"
)

func TestIsValidPoint(t *testing.T) {
	// Test data with known coordinates
	baseTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		current  *model.Log
		next     *model.Log
		expected bool
	}{
		{
			name:     "nil current point",
			current:  nil,
			next:     &model.Log{Latitud: 40.7128, Longitud: -74.0060, Date: baseTime},
			expected: false,
		},
		{
			name:     "nil next point",
			current:  &model.Log{Latitud: 40.7128, Longitud: -74.0060, Date: baseTime},
			next:     nil,
			expected: false,
		},
		{
			name:     "current date after next date",
			current:  &model.Log{Latitud: 40.7128, Longitud: -74.0060, Date: baseTime.Add(time.Hour)},
			next:     &model.Log{Latitud: 40.7589, Longitud: -73.9851, Date: baseTime},
			expected: false,
		},
		{
			name:     "zero time difference",
			current:  &model.Log{Latitud: 40.7128, Longitud: -74.0060, Date: baseTime},
			next:     &model.Log{Latitud: 40.7589, Longitud: -73.9851, Date: baseTime},
			expected: false,
		},
		{
			name:     "same coordinates (zero distance)",
			current:  &model.Log{Latitud: 40.7128, Longitud: -74.0060, Date: baseTime},
			next:     &model.Log{Latitud: 40.7128, Longitud: -74.0060, Date: baseTime.Add(time.Hour)},
			expected: false,
		},
		{
			name:     "valid speed - normal movement",
			current:  &model.Log{Latitud: 40.7128, Longitud: -74.0060, Date: baseTime},                // NYC (approx)
			next:     &model.Log{Latitud: 40.7589, Longitud: -73.9851, Date: baseTime.Add(time.Hour)}, // Central Park (approx)
			expected: true,                                                                            // Distance ~5.5km in 1 hour = ~5.5 km/h (well below threshold)
		},
		{
			name:     "valid speed - highway speed",
			current:  &model.Log{Latitud: 40.7128, Longitud: -74.0060, Date: baseTime},                // NYC
			next:     &model.Log{Latitud: 41.0128, Longitud: -74.0060, Date: baseTime.Add(time.Hour)}, // ~33km north
			expected: true,                                                                            // Distance ~33km in 1 hour = ~33 km/h (below threshold)
		},
		{
			name:     "invalid speed - too fast",
			current:  &model.Log{Latitud: 40.7128, Longitud: -74.0060, Date: baseTime},                // NYC
			next:     &model.Log{Latitud: 44.0128, Longitud: -74.0060, Date: baseTime.Add(time.Hour)}, // ~366km north
			expected: false,                                                                           // Distance ~366km in 1 hour = ~366 km/h (above 350 km/h threshold)
		}, {
			name:     "edge case - just below threshold",
			current:  &model.Log{Latitud: 40.7128, Longitud: -74.0060, Date: baseTime},                // NYC
			next:     &model.Log{Latitud: 43.8000, Longitud: -74.0060, Date: baseTime.Add(time.Hour)}, // ~345km north
			expected: true,                                                                            // Distance ~345km in 1 hour = ~345 km/h (just below threshold)
		},
		{
			name:     "valid speed - short time interval",
			current:  &model.Log{Latitud: 40.7128, Longitud: -74.0060, Date: baseTime},
			next:     &model.Log{Latitud: 40.7138, Longitud: -74.0050, Date: baseTime.Add(time.Minute)}, // Very small movement
			expected: true,                                                                              // Small distance in short time = reasonable speed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidPoint(tt.current, tt.next)
			if result != tt.expected {
				t.Errorf("IsValidPoint() = %v, expected %v", result, tt.expected)

				// Debug information for failed tests
				if tt.current != nil && tt.next != nil {
					timeDiff := tt.next.Date.Sub(tt.current.Date).Hours()
					t.Logf("Time difference: %.6f hours", timeDiff)

					if timeDiff > 0 {
						// Calculate what the function would compute
						currentPoint := s2.LatLngFromDegrees(tt.current.Latitud, tt.current.Longitud)
						nextPoint := s2.LatLngFromDegrees(tt.next.Latitud, tt.next.Longitud)
						distance := nextPoint.Distance(currentPoint).Radians() * earthRadiusKm
						speed := distance / timeDiff
						t.Logf("Distance: %.6f km, Speed: %.6f km/h", distance, speed)
					}
				}
			}
		})
	}
}
