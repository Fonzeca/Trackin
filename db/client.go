package db

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//Sirve para obener el objeto para interactuar con la base de datos
func ObtenerConexionDb() (*gorm.DB, func() error, error) {
	host := os.Getenv("trackin-db-host")
	if host == "" {
		host = "vps-1791261-x.dattaweb.com:3306"
	}

	user := os.Getenv("trackin-db-user")
	if user == "" {
		user = "root"
	}
	pass := os.Getenv("trackin-db-pass")
	if pass == "" {
		pass = "almacen.C12"
	}

	dsn := user + ":" + pass + "@tcp(" + host + ")/trackin?parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	sqlDb, _ := db.DB()
	return db, sqlDb.Close, nil
}
