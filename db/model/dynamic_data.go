package model

import "time"

type LastLogView struct {
	Imei     string    `json:"imei"`
	Latitutd float64   `json:"latitud"`
	Longitud float64   `json:"longitud"`
	Speed    float32   `json:"speed"`
	Date     time.Time `json:"date"`
}
