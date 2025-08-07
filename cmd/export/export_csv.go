package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	manager "github.com/Fonzeca/Trackin/internal/core/managers"
	"github.com/Fonzeca/Trackin/internal/infrastructure/database"
	"github.com/Fonzeca/Trackin/internal/infrastructure/database/model"
	"github.com/Fonzeca/Trackin/internal/infrastructure/geolocation"
	"github.com/golang/geo/s2"
)

// interface
type PointIntersection struct {
	log   *model.Log
	zones []*model.Zona
}

func main() {
	// Define command line flags
	var (
		startDate = flag.String("start", "", "Start date (format: YYYY-MM-DD HH:MM:SS, ex: 2025-12-25 14:30:00)")
		endDate   = flag.String("end", "", "End date (format: YYYY-MM-DD HH:MM:SS, ex: 2025-12-31 23:59:59)")
		imei      = flag.String("imei", "", "Device IMEI")
		output    = flag.String("output", "route.csv", "CSV output file")
	)

	flag.Parse()

	// Validate required parameters
	if *startDate == "" || *endDate == "" || *imei == "" {
		fmt.Println("Usage: export_csv -start='YYYY-MM-DD HH:MM:SS' -end='YYYY-MM-DD HH:MM:SS' -imei='123456789' [-output='route.csv']")
		fmt.Println("Example: export_csv -start='2025-12-15 08:00:00' -end='2025-12-31 18:30:00' -imei='356307042441013'")
		flag.Usage()
		os.Exit(1)
	}

	// Validate date formats
	startDateTime, err := time.Parse("2006-01-02 15:04:05", *startDate)
	if err != nil {
		log.Fatalf("Error in start date format: %v", err)
	}

	endDateTime, err := time.Parse("2006-01-02 15:04:05", *endDate)
	if err != nil {
		log.Fatalf("Error in end date format: %v", err)
	}

	// Validate that start date is before end date
	if startDateTime.After(endDateTime) {
		log.Fatal("Start date must be before end date")
	}

	// Initialize database
	database.InitDB()
	defer database.CloseDB()

	// Get GPS logs for IMEI in date range
	logs, err := getGPSLogs(*imei, startDateTime, endDateTime)
	if err != nil {
		log.Fatalf("Error getting GPS logs: %v", err)
	}

	if len(logs) == 0 {
		fmt.Printf("No GPS data found for IMEI %s between %s and %s\n", *imei, *startDate, *endDate)
		return
	}

	// Get zones
	zonesManager := manager.GetManagerContainer().GetZonasManager()
	zones, err := zonesManager.GetZonesByEmpresaId(971) // Assuming "971" is the company ID

	if err != nil {
		log.Fatalf("Error getting zones: %v", err)
	}

	if len(zones) == 0 {
		fmt.Println("No zones found for company ID 971")
		return
	}

	// Map zone ID to zone request
	zoneMap := make(map[int32]*model.Zona)
	for _, zone := range zones {
		zoneMap[zone.ID] = &zone
	}

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

	// Intersect logs with zones
	intersections, err := intersectLogsAndZones(logs, zoneMap, loopMap)
	if err != nil {
		log.Fatalf("Error intersecting logs with zones: %v", err)
	}

	// Export to CSV
	err = exportToCSV(intersections, *output)
	if err != nil {
		log.Fatalf("Error exporting to CSV: %v", err)
	}

	fmt.Printf("Export successful: %d records exported to %s\n", len(intersections), *output)
}

// getGPSLogs gets GPS logs for the IMEI in the specified date range
func getGPSLogs(imei string, startDate, endDate time.Time) ([]model.Log, error) {

	log.Println("Getting GPS logs for IMEI:", imei, "from", startDate, "to", endDate)

	routesManager := manager.GetManagerContainer().GetRoutesManager()

	routes, err := routesManager.GetRouteByImei(model.RouteRequest{
		Imei: imei,
		From: startDate.String(),
		To:   endDate.String(),
	})

	if err != nil {
		return nil, fmt.Errorf("error getting route by IMEI: %w", err)
	}

	log.Println("Routes obtained:", len(routes), "records")

	// Filter valid GPS logs
	if len(routes) == 0 {
		return nil, nil
	}

	// Log routes as JSON
	routesJSON, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		log.Printf("Error serializing routes to JSON: %v", err)
	}

	// Export JSON to file
	file, err := os.Create("routes.json")
	if err != nil {
		log.Printf("Error creating routes.json file: %v", err)
		return nil, err
	}
	defer file.Close()

	// Write JSON to file
	if _, err := file.Write(routesJSON); err != nil {
		log.Printf("Error writing to routes.json file: %v", err)
		return nil, err
	}

	logs := make([]model.Log, 0)

	for _, route := range routes {

		if route.Type == "Parada" {

			// Combine date and time into a single string
			dateTimeStr := fmt.Sprintf("%s %s", route.FromDate, route.FromHour)
			// Parse the combined string to time.Time
			dateTime, err := time.Parse("2006-01-02 15:04", dateTimeStr)
			if err != nil {
				log.Printf("Error parsing date and time: %v", err)
				continue
			}
			logs = append(logs, model.Log{
				Imei:     imei,
				Latitud:  route.Latitud,
				Longitud: route.Longitud,
				Date:     dateTime,
				Speed:    0,
			})

		} else if route.Type == "Viaje" {

			for _, data := range route.Data {

				// Convert epoch ms to time.Time in UTC
				dateTime := time.UnixMilli(data.Timestamp).UTC()

				logs = append(logs, model.Log{
					Imei:     imei,
					Latitud:  data.Latitud,
					Longitud: data.Longitud,
					Date:     dateTime,
					Speed:    data.Speed,
				})
			}

		}

	}

	log.Println("GPS logs obtained:", len(logs), "records")

	return logs, nil
}

// getZoneNames extracts zone names and concatenates them with semicolon
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

// exportToCSV exports the intersections to a CSV file
func exportToCSV(intersections []PointIntersection, filename string) error {

	log.Println("Exporting intersections to CSV:", filename)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	headers := []string{
		"IMEI",
		"Fecha",
		"Latitud",
		"Longitud",
		"Velocidad (km/h)",
		"Zonas",
		"Velocidad Maxima",
		"Exceso de Velocidad",
		"Link maps",
	}

	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write data
	for _, intersection := range intersections {
		log := intersection.log

		// Get zone names using helper function
		zonesStr := getZoneNames(intersection.zones)

		// Get the lowest speed limit using helper function
		minSpeedLimit := getMinSpeedLimit(intersection.zones)

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
			fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%f,%f", log.Latitud, log.Longitud),
		}

		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func intersectLogsAndZones(logs []model.Log, zoneMap map[int32]*model.Zona, loopMap map[int32]*s2.Loop) ([]PointIntersection, error) {
	if len(logs) == 0 || len(zoneMap) == 0 || len(loopMap) == 0 {
		return nil, fmt.Errorf("logs, zones or loops cannot be empty")
	}

	var intersections []PointIntersection

	for _, log := range logs {
		intersec := PointIntersection{
			log:   &log,
			zones: make([]*model.Zona, 0),
		}

		for zoneId, loop := range loopMap {
			if loop.ContainsPoint(s2.PointFromLatLng(s2.LatLngFromDegrees(log.Latitud, log.Longitud))) {
				if zone, exists := zoneMap[zoneId]; exists {
					intersec.zones = append(intersec.zones, zone)
				}
			}
		}

		// Add intersection to array (even if it has no zones)
		intersections = append(intersections, intersec)
	}

	return intersections, nil
}
