package json

import "time"

type SimplyData struct {
	Imei         string    `json:"imei"`
	ProtocolType int32     `json:"protocolHeadType"`
	Longitude    float64   `json:"longitude"`
	Latitude     float64   `json:"latitude"`
	Date         time.Time `json:"date"`
	Speed        float32   `json:"speed"`
	AnalogInput1 float32   `json:"analogInput1"`
	DeviceTemp   int32     `json:"deviceTemp"`
	Mileage      int32     `json:"mileage"`
	GpsWorking   bool      `json:"gpsWorking"`
	IsHistory    bool      `json:"isHistoryData"`
	EngineStatus bool      `json:"iopIgnition"`
	Azimuth      int32     `json:"azimuth"`
	PayLoad      string
	CloseChannel bool
}

type SimplyDataLocation struct {
	Imei        string  `protobuf:"bytes,1,opt,name=imei,proto3" json:"imei,omitempty"`
	Latitude    float64 `protobuf:"fixed64,2,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude   float64 `protobuf:"fixed64,3,opt,name=longitude,proto3" json:"longitude,omitempty"`
	TimestampMs int64   `protobuf:"varint,4,opt,name=timestamp_ms,json=timestampMs,proto3" json:"timestamp_ms,omitempty"`
	Speed       int32   `protobuf:"varint,5,opt,name=speed,proto3" json:"speed,omitempty"`
	EngineOn    bool    `protobuf:"varint,6,opt,name=engine_on,json=engineOn,proto3" json:"engine_on,omitempty"`
	Azimuth     int32   `protobuf:"varint,7,opt,name=azimuth,proto3" json:"azimuth,omitempty"`
	Payload     string  `protobuf:"bytes,99,opt,name=payload,proto3" json:"payload,omitempty"`
}
