package model

import "time"

type LastLogView struct {
	Imei     string    `json:"imei"`
	Latitutd float64   `json:"latitud"`
	Longitud float64   `json:"longitud"`
	Speed    float32   `json:"speed"`
	Date     time.Time `json:"date"`
}

type RouteView struct {
	Imei string          `json:"imei"`
	From string          `json:"from"`
	To   string          `json:"to"`
	Data []RouteDataView `json:"data"`
}

type RouteDataView struct {
	Location
	Date  time.Time `json:"date"`
	Speed float32   `json:"speed"`
}

type Location struct {
	Latitutd float64 `json:"latitud"`
	Longitud float64 `json:"longitud"`
}
