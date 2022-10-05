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
	EngineStatus bool  `json:"engine_status"`
	Azimuth      int32 `json:"azimuth,omitempty"`
}

type RouteView struct {
	Id       int32  `json:"id"`
	Type     string `json:"type"`
	FromDate string `json:"fromDate"`
	ToDate   string `json:"toDate"`
	FromHour string `json:"fromHour"`
	ToHour   string `json:"toHour"`
}

type RouteDataView struct {
	Location
	Speed   float32 `json:"speed"`
	Azimuth int32   `json:"azimuth,omitempty"`
}

type StopView struct {
	RouteView
	Location
}

type MoveView struct {
	RouteView
	KM   int32           `json:"km"`
	Data []RouteDataView `json:"data"`
}

type RouteRequest struct {
	Imei string `json:"imei"`
	From string `json:"from"`
	To   string `json:"to"`
}

type ImeisBody struct {
	Imeis []string `json:"imeis"`
}

type ZoneRequest struct {
	EmpresaId     int32  `json:"empresa_id,omitempty"`
	ColorLinea    string `json:"color_linea,omitempty"`
	ColorRelleno  string `json:"color_relleno,omitempty"`
	Puntos        string `json:"puntos,omitempty"`
	Nombre        string `json:"nombre,omitempty"`
	VehiculosIds  []int  `json:"vehiculos_ids"`
	AvisarEntrada bool   `json:"avisar_entrada,omitempty"`
	AvisarSalida  bool   `json:"avisar_salida,omitempty"`
}

type ZoneVehiclesView struct {
	ZonaID        bool    `json:"zona_id,omitempty"`
	VehiculosId   []int32 `json:"vehiculos_ids,omitempty"`
	AvisarEntrada bool    `json:"avisar_entrada,omitempty"`
	AvisarSalida  bool    `json:"avisar_salida,omitempty"`
}
