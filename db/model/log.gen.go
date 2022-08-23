// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameLog = "log"

// Log mapped from table <log>
type Log struct {
	ID           int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Imei         string    `gorm:"column:imei;not null" json:"imei"`
	ProtocolType int32     `gorm:"column:protocol_type;not null" json:"protocol_type"`
	Latitud      *float64   `gorm:"column:latitud" json:"latitud"`
	Longitud     *float64   `gorm:"column:longitud" json:"longitud"`
	Date         time.Time `gorm:"column:date;not null" json:"date"`
	Speed        float32   `gorm:"column:speed" json:"speed"`
	DeviceTemp   int32     `gorm:"column:device_temp" json:"device_temp"`
	Mileage      int32     `gorm:"column:mileage" json:"mileage"`
	IsGps        bool      `gorm:"column:is_gps" json:"is_gps"`
	IsHistory    bool      `gorm:"column:is_history" json:"is_history"`
	EngineStatus *bool      `gorm:"column:engine_status" json:"engine_status"`
	Azimuth      *int32     `gorm:"column:azimuth" json:"azimuth"`
}

// TableName Log's table name
func (*Log) TableName() string {
	return TableNameLog
}
