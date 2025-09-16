package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Fonzeca/Trackin/internal/infrastructure/database"
	"github.com/Fonzeca/Trackin/internal/infrastructure/database/model"
	"github.com/Fonzeca/Trackin/internal/infrastructure/geolocation"
)

type LogValidationResult struct {
	Log         *model.Log
	IsValid     bool
	PreviousLog *model.Log
}

// Estructuras para el JSON del mapa
type MapPoint struct {
	ID    string  `json:"id"`
	Label string  `json:"label"`
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
}

type MapLine struct {
	ID     string     `json:"id"`
	Label  string     `json:"label"`
	Color  string     `json:"color"`
	Points []MapPoint `json:"points"`
}

type MapData struct {
	MapLines []MapLine `json:"mapLines"`
}

func main() {
	// Definir flags de línea de comandos
	var (
		startDate = flag.String("start", "", "Fecha de inicio (formato: YYYY-MM-DD HH:MM:SS, ej: 2025-12-25 14:30:00)")
		endDate   = flag.String("end", "", "Fecha de fin (formato: YYYY-MM-DD HH:MM:SS, ej: 2025-12-31 23:59:59)")
		imei      = flag.String("imei", "", "IMEI del dispositivo")
		dryRun    = flag.Bool("dry-run", true, "Solo mostrar resultados sin eliminar registros (por defecto: true)")
		verbose   = flag.Bool("verbose", false, "Mostrar información detallada")
	)

	flag.Parse()

	// Validar parámetros requeridos
	if *startDate == "" || *endDate == "" || *imei == "" {
		fmt.Println("Uso: debug_logs -start='YYYY-MM-DD HH:MM:SS' -end='YYYY-MM-DD HH:MM:SS' -imei='123456789' [-dry-run=true] [-verbose] [-map-url='URL']")
		fmt.Println("Ejemplo: debug_logs -start='2025-12-15 08:00:00' -end='2025-12-31 18:30:00' -imei='356307042441013' -verbose -map-url='https://mi-mapa.com'")
		fmt.Println("\nValidación: Usa IsValidPointFlexible (solo verifica velocidades excesivas > 500 km/h)")
		fmt.Println("\nFlags:")
		fmt.Println("  -dry-run=false : Eliminar registros inválidos de la base de datos")
		fmt.Println("  -verbose       : Mostrar información detallada de cada validación")
		fmt.Println("  -map-url       : URL base para visualización en mapa (por defecto: https://localhost:3000)")
		flag.Usage()
		os.Exit(1)
	}

	// Validar formatos de fecha
	startDateTime, err := time.Parse("2006-01-02 15:04:05", *startDate)
	if err != nil {
		log.Fatalf("Error en formato de fecha de inicio: %v", err)
	}

	endDateTime, err := time.Parse("2006-01-02 15:04:05", *endDate)
	if err != nil {
		log.Fatalf("Error en formato de fecha de fin: %v", err)
	}

	if startDateTime.After(endDateTime) {
		log.Fatalf("La fecha de inicio debe ser anterior a la fecha de fin")
	}

	// Inicializar base de datos
	database.InitDB()
	defer database.CloseDB()

	// Obtener logs del IMEI en el rango de fechas usando query RAW
	var logs []*model.Log
	rawSQL := `
		SELECT id, imei, protocol_type, latitud, longitud, date, speed, 
		       analog_input_1, device_temp, mileage, is_gps, is_history, 
		       engine_status, azimuth, payload 
		FROM log 
		WHERE imei = ? AND date >= ? AND date <= ? 
		ORDER BY date ASC
	`

	err = database.DB.Raw(rawSQL, *imei, startDateTime, endDateTime).Find(&logs).Error

	if err != nil {
		log.Fatalf("Error al obtener logs: %v", err)
	}

	if len(logs) == 0 {
		fmt.Printf("No se encontraron logs para el IMEI %s en el rango de fechas especificado\n", *imei)
		os.Exit(0)
	}

	fmt.Printf("📊 Analizando %d logs para IMEI: %s\n", len(logs), *imei)
	fmt.Printf("📅 Rango: %s a %s\n", startDateTime.Format("2006-01-02 15:04:05"), endDateTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("🔍 Modo: %s\n", func() string {
		if *dryRun {
			return "DRY-RUN (solo análisis)"
		}
		return "ELIMINACIÓN ACTIVA"
	}())
	fmt.Printf("🚀 Validación: FLEXIBLE (solo velocidades > 500 km/h se consideran inválidas)\n")
	fmt.Println("=" + string(make([]rune, 80)))

	// Validar logs
	results := validateLogs(logs, *verbose)

	// Mostrar resumen
	showSummary(results, logs)

	// Eliminar logs inválidos si no es dry-run
	if !*dryRun {
		deleteInvalidLogs(results)
	}
}

func validateLogs(logs []*model.Log, verbose bool) []LogValidationResult {
	var results []LogValidationResult
	validCount := 0
	invalidCount := 0

	for i, currentLog := range logs {
		result := LogValidationResult{
			Log:     currentLog,
			IsValid: true,
		}

		// Primer log siempre es válido (no hay anterior para comparar)
		if i == 0 {
			results = append(results, result)
			validCount++
			continue
		}

		previousLog := logs[i-1]
		result.PreviousLog = previousLog

		// Validar usando IsValidPointFlexible
		isValid := geolocation.IsValidPointFlexible(previousLog, currentLog)

		if !isValid {
			result.IsValid = false
			invalidCount++

			if verbose {
				fmt.Printf("❌ Log #%d [%s]: INVÁLIDO\n",
					i+1, currentLog.Date.Format("2006-01-02 15:04:05"))
				fmt.Printf("    📍 Lat: %.6f, Lon: %.6f\n", currentLog.Latitud, currentLog.Longitud)
			}
		} else {
			validCount++
		}

		results = append(results, result)
	}

	return results
}

func showSummary(results []LogValidationResult, allLogs []*model.Log) {
	validCount := 0
	invalidCount := 0

	for _, result := range results {
		if result.IsValid {
			validCount++
		} else {
			invalidCount++
		}
	}

	fmt.Println("\n📋 RESUMEN DE VALIDACIÓN")
	fmt.Println("=" + string(make([]rune, 50)))
	fmt.Printf("📈 Logs totales:     %d\n", len(results))
	fmt.Printf("✅ Logs válidos:     %d (%.1f%%)\n", validCount, float64(validCount)/float64(len(results))*100)
	fmt.Printf("❌ Logs inválidos:   %d (%.1f%%)\n", invalidCount, float64(invalidCount)/float64(len(results))*100)

	if invalidCount > 0 {
		fmt.Println("\n🗑️  LOGS INVÁLIDOS ENCONTRADOS:")
		var mapLines []MapLine

		for i, result := range results {
			if !result.IsValid {
				fmt.Printf("  • Log #%d [%s]: ID=%d\n",
					i+1, result.Log.Date.Format("2006-01-02 15:04:05"), result.Log.ID)

				// Crear ruta con contexto (2 puntos antes + punto inválido + 2 puntos después)
				mapLine := createMapLineForInvalidPoint(allLogs, i, result.Log.ID)
				if mapLine != nil {
					mapLines = append(mapLines, *mapLine)
				}
			}
		}

		// Generar JSON del mapa si hay puntos inválidos
		if len(mapLines) > 0 {
			mapData := MapData{MapLines: mapLines}
			jsonData, err := json.MarshalIndent(mapData, "", "  ")
			if err != nil {
				log.Printf("Error al generar JSON: %v", err)
			} else {
				fmt.Printf("\n🗺️  JSON PARA VISUALIZAR PUNTOS INVÁLIDOS EN EL MAPA:\n")
				fmt.Printf("Copia este JSON en tu aplicación de mapas:\n")
				fmt.Printf("\n%s\n", string(jsonData))
			}
		}
	}
}

func createMapLineForInvalidPoint(allLogs []*model.Log, invalidIndex int, invalidLogID int32) *MapLine {
	if invalidIndex < 0 || invalidIndex >= len(allLogs) {
		return nil
	}

	var points []MapPoint
	colors := []string{"#ef4444", "#f97316", "#eab308", "#22c55e", "#3b82f6", "#8b5cf6"}

	// Calcular rango de puntos (2 antes + punto inválido + 2 después)
	startIdx := invalidIndex - 2
	if startIdx < 0 {
		startIdx = 0
	}

	endIdx := invalidIndex + 2
	if endIdx >= len(allLogs) {
		endIdx = len(allLogs) - 1
	}

	// Crear puntos para la ruta
	for i := startIdx; i <= endIdx; i++ {
		log := allLogs[i]
		var label string

		if log.ID == invalidLogID {
			label = "🚨 PUNTO INVÁLIDO"
		} else if i < invalidIndex {
			label = fmt.Sprintf("📍 Anterior -%d", invalidIndex-i)
		} else {
			label = fmt.Sprintf("📍 Posterior +%d", i-invalidIndex)
		}

		point := MapPoint{
			ID:    fmt.Sprintf("%d", log.ID),
			Label: fmt.Sprintf("%s [%s]", label, log.Date.Format("15:04:05")),
			Lat:   log.Latitud,
			Lng:   log.Longitud,
		}
		points = append(points, point)
	}

	// Seleccionar color basado en el ID del log inválido
	colorIndex := int(invalidLogID) % len(colors)

	mapLine := &MapLine{
		ID:     fmt.Sprintf("ruta_invalida_%d", invalidLogID),
		Label:  fmt.Sprintf("Punto Inválido ID:%d [%s]", invalidLogID, allLogs[invalidIndex].Date.Format("2006-01-02 15:04:05")),
		Color:  colors[colorIndex],
		Points: points,
	}

	return mapLine
}

func deleteInvalidLogs(results []LogValidationResult) {
	var idsToDelete []int32

	for _, result := range results {
		if !result.IsValid {
			idsToDelete = append(idsToDelete, result.Log.ID)
		}
	}

	if len(idsToDelete) == 0 {
		fmt.Println("\n✨ No hay logs inválidos para eliminar")
		return
	}

	fmt.Printf("\n🗑️  Eliminando %d logs inválidos...\n", len(idsToDelete))

	// Usar query RAW para eliminar
	rawSQL := "DELETE FROM log WHERE id IN (?)"
	result := database.DB.Exec(rawSQL, idsToDelete)

	if result.Error != nil {
		log.Fatalf("Error al eliminar logs: %v", result.Error)
	}

	fmt.Printf("✅ Se eliminaron %d registros exitosamente\n", result.RowsAffected)
}
