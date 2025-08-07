package model

import "time"

type Location struct {
	Latitutd float64 `json:"latitud,omitempty"`
	Longitud float64 `json:"longitud,omitempty"`
}

type LastLogView struct {
	Imei string `json:"imei"`
	Location
	Speed float32   `json:"speed"`
	Date  time.Time `json:"date"`
}

type StateLogView struct {
	Imei string `json:"imei"`
	Location
	EngineStatus bool      `json:"engine_status"`
	Azimuth      int32     `json:"azimuth,omitempty"`
	Date         time.Time `json:"date"`
}

type RouteRequest struct {
	Imei     string  `json:"imei"`
	From     string  `json:"from"`
	To       string  `json:"to"`
	ZonesIds []int32 `json:"zones_ids,omitempty"`
}

type ImeisBody struct {
	Imeis []string `json:"imeis"`
}

type ZoneView struct {
	Id              int32   `json:"id,omitempty"`
	EmpresaId       int32   `json:"empresa_id,omitempty"`
	ColorLinea      string  `json:"color_linea,omitempty"`
	ColorRelleno    string  `json:"color_relleno,omitempty"`
	Puntos          string  `json:"puntos,omitempty"`
	Nombre          string  `json:"nombre,omitempty"`
	Imei            string  `json:"imei"`
	AvisarEntrada   bool    `json:"avisar_entrada,omitempty"`
	AvisarSalida    bool    `json:"avisar_salida,omitempty"`
	VelocidadMaxima float64 `json:"velocidad_maxima,omitempty"`
}

type ZoneRequest struct {
	Id              int32    `json:"id,omitempty"`
	EmpresaId       int32    `json:"empresa_id,omitempty"`
	ColorLinea      string   `json:"color_linea,omitempty"`
	ColorRelleno    string   `json:"color_relleno,omitempty"`
	Puntos          string   `json:"puntos,omitempty"`
	Nombre          string   `json:"nombre,omitempty"`
	Imeis           []string `json:"imeis"`
	AvisarEntrada   bool     `json:"avisar_entrada,omitempty"`
	AvisarSalida    bool     `json:"avisar_salida,omitempty"`
	VelocidadMaxima float64  `json:"velocidad_maxima,omitempty"`
}

type ZoneVehiclesView struct {
	ZonaID        bool    `json:"zona_id,omitempty"`
	VehiculosId   []int32 `json:"vehiculos_ids,omitempty"`
	AvisarEntrada bool    `json:"avisar_entrada,omitempty"`
	AvisarSalida  bool    `json:"avisar_salida,omitempty"`
}

type ZoneNotification struct {
	Imei      string `json:"imei,omitempty"`
	ZoneName  string `json:"zone_name,omitempty"`
	ZoneID    int    `json:"zone_id,omitempty"`
	EventType string `json:"event_type,omitempty"`
}

type GpsPoint struct {
	Azimuth   int32   `json:"azimuth"`
	Latitud   float64 `json:"latitud"`
	Longitud  float64 `json:"longitud"`
	Speed     float32 `json:"speed"`
	Timestamp int64   `json:"timestamp"`
}

type GpsRouteData struct {
	Id       int32      `json:"id"`
	Type     string     `json:"type"`
	FromDate string     `json:"fromDate"`
	ToDate   string     `json:"toDate"`
	FromHour string     `json:"fromHour"`
	ToHour   string     `json:"toHour"`
	Duration string     `json:"duration"`
	Latitud  float64    `json:"latitud,omitempty"`
	Longitud float64    `json:"longitud,omitempty"`
	Km       int32      `json:"km,omitempty"`
	Data     []GpsPoint `json:"data,omitempty"`
}

type SummaryRequest struct {
	Imei      string `json:"imei"`
	FromDate  int64  `json:"fromDate"`
	ToDate    int64  `json:"toDate"`
	EmpresaId int32  `json:"empresa_id,omitempty"`
}

type PointIntersection struct {
	Log   *Log
	Zones []*Zona
}
