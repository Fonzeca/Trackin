package json

import "time"

type SimplyData struct {
	Imei      string    `json:"imei"`
	Longitude float64   `json:"longitude"`
	Latitude  float64   `json:"latitude"`
	Date      time.Time `json:"date"`
}
