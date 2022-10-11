package db

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Sirve para obener el objeto para interactuar con la base de datos
func ObtenerConexionDb() (*gorm.DB, func() error, error) {
	//Cambiarlo por Viper
	host := os.Getenv("trackinDbHost")
	if host == "" {
		host = "vps-1791261-x.dattaweb.com:3306"
	}

	user := os.Getenv("trackinDbUser")
	if user == "" {
		user = "root"
	}
	pass := os.Getenv("trackinDbPass")
	if pass == "" {
		pass = "almacen.C12"
	}

	dsn := user + ":" + pass + "@tcp(" + host + ")/trackin?parseTime=false"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	sqlDb, _ := db.DB()
	return db, sqlDb.Close, nil
}
