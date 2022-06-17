package json

import "time"

type SimplyData struct {
	Imei         string    `json:"imei"`
	ProtocolType int32     `json:"protocolHeadType"`
	Longitude    float64   `json:"longitude"`
	Latitude     float64   `json:"latitude"`
	Date         time.Time `json:"date"`
	Speed        float32   `json:"speed"`
	DeviceTemp   int32     `json:"deviceTemp"`
	Mileage      int32     `json:"mileage"`
	GpsWorking   bool      `json:"gpsWorking"`
	IsHistory    bool      `json:"isHistoryData"`
	EngineStatus bool      `json:"iopIgnition"`
}
