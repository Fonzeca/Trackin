package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	//TODO: Cambiarlo por Viper
	host := os.Getenv("trackinDbHost")
	if host == "" {
		host = "vps-2367826-x.dattaweb.com:3306"
	}

	user := os.Getenv("trackinDbUser")
	if user == "" {
		user = "root"
	}
	pass := os.Getenv("trackinDbPass")
	if pass == "" {
		pass = "carmind-db"
	}

	dsn := user + ":" + pass + "@tcp(" + host + ")/trackin?parseTime=True"

	var err error
	dbMinus, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	DB = dbMinus

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error al obtener el *sql.DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Testeamos la conexión
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Error haciendo ping a la base: %v", err)
	}

	log.Println("Conexión a la base exitosa")
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err == nil {
		sqlDB.Close()
	}
}
