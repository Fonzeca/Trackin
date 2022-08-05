package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//Sirve para obener el objeto para interactuar con la base de datos
func ObtenerConexionDb() (*gorm.DB, error) {
	dsn := "root:almacen.C12@tcp(vps-1791261-x.dattaweb.com:3306)/trackin?charset=utf8&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
