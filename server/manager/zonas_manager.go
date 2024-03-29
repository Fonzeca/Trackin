package manager

import (
	"strconv"

	"github.com/Fonzeca/Trackin/db"
	"github.com/Fonzeca/Trackin/db/model"
	"gorm.io/gorm"
)

type IZonasManager interface {
	GetZonesByEmpresaId(idParam string) ([]model.ZoneRequest, error)
	CreateZone(zoneRequest model.ZoneRequest) error
	EditZoneById(idParam string, zoneRequest model.ZoneRequest) error
	DeleteZoneById(idParam string) error
	GetZoneConfigByImei(imei string) ([]model.ZoneView, error)
}

var ZonasManager IZonasManager = newZonasManager()

type zonasManager struct {
}

func newZonasManager() IZonasManager {
	return &zonasManager{}
}

func (ma *zonasManager) GetZonesByEmpresaId(idParam string) ([]model.ZoneRequest, error) {
	db, close, err := db.ObtenerConexionDb()
	defer close()

	if err != nil {
		return nil, err
	}

	//Scamos el id de la empresa
	id, idParseErr := strconv.Atoi(idParam)

	if idParseErr != nil {
		return []model.ZoneRequest{}, idParseErr
	}

	zones := []model.ZoneView{}

	tx := db.Model(&model.Zona{}).Select("zona.id, zona.empresa_id, zona.color_linea, zona.color_relleno, zona.puntos, zona.nombre, zona_vehiculos.imei, zona_vehiculos.avisar_entrada, zona_vehiculos.avisar_salida").Joins("join zona_vehiculos on zona.id = zona_vehiculos.zona_id").Where("empresa_id = ?", id).Order("zona.id desc").Scan(&zones)

	if tx.Error != nil {
		return nil, tx.Error
	}

	zonesWithVehicles := []model.ZoneRequest{}

	var previousZoneId int32 = -1
	for i, zone := range zones {
		if zone.Id != previousZoneId || i == 0 {
			zonesWithVehicles = append(zonesWithVehicles, model.ZoneRequest{
				Id:            zone.Id,
				EmpresaId:     zone.EmpresaId,
				ColorLinea:    zone.ColorLinea,
				ColorRelleno:  zone.ColorRelleno,
				Puntos:        zone.Puntos,
				Nombre:        zone.Nombre,
				Imeis:         []string{zone.Imei},
				AvisarEntrada: zone.AvisarEntrada,
				AvisarSalida:  zone.AvisarSalida,
			})
		} else {
			zonesWithVehicles[len(zonesWithVehicles)-1].Imeis = append(zonesWithVehicles[len(zonesWithVehicles)-1].Imeis, zone.Imei)
		}
		previousZoneId = zone.Id
	}

	zones = []model.ZoneView{}

	tx = db.Model(&model.Zona{}).Joins("left outer join zona_vehiculos on zona.id = zona_vehiculos.zona_id").Where("empresa_id = ? AND zona_vehiculos.zona_id is null", id).Order("id desc").Scan(&zones)

	zonesWithoutVehciles := []model.ZoneRequest{}

	for _, zone := range zones {
		zonesWithoutVehciles = append(zonesWithoutVehciles, model.ZoneRequest{
			Id:            zone.Id,
			EmpresaId:     zone.EmpresaId,
			ColorLinea:    zone.ColorLinea,
			ColorRelleno:  zone.ColorRelleno,
			Puntos:        zone.Puntos,
			Nombre:        zone.Nombre,
			Imeis:         []string{zone.Imei},
			AvisarEntrada: zone.AvisarEntrada,
			AvisarSalida:  zone.AvisarSalida,
		})
	}
	zonesWithVehicles = append(zonesWithVehicles, zonesWithoutVehciles...)

	return zonesWithVehicles, tx.Error
}

func (ma *zonasManager) CreateZone(zoneRequest model.ZoneRequest) error {
	db, close, err := db.ObtenerConexionDb()
	defer close()

	if err != nil {
		return err
	}

	transactionErr := db.Transaction(func(tx *gorm.DB) error {

		zone := model.Zona{
			EmpresaID:    int32(zoneRequest.EmpresaId),
			ColorLinea:   zoneRequest.ColorLinea,
			ColorRelleno: zoneRequest.ColorRelleno,
			Puntos:       zoneRequest.Puntos,
			Nombre:       zoneRequest.Nombre,
		}

		if err := tx.Create(&zone).Error; err != nil {
			return err
		}

		if len(zoneRequest.Imeis) > 0 {
			zonesWithVehicles := []model.ZonaVehiculo{}
			for _, imei := range zoneRequest.Imeis {
				zonesWithVehicles = append(zonesWithVehicles, model.ZonaVehiculo{
					ZonaID:        zone.ID,
					Imei:          imei,
					AvisarEntrada: zoneRequest.AvisarEntrada,
					AvisarSalida:  zoneRequest.AvisarSalida,
				})
			}

			if err := tx.Create(&zonesWithVehicles).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return transactionErr
}

func (ma *zonasManager) EditZoneById(idParam string, zoneRequest model.ZoneRequest) error {
	db, close, err := db.ObtenerConexionDb()
	defer close()

	if err != nil {
		return err
	}

	id, idParseErr := strconv.Atoi(idParam)

	if idParseErr != nil {
		return idParseErr
	}

	zone := model.Zona{
		ID:           int32(id),
		EmpresaID:    int32(zoneRequest.EmpresaId),
		ColorLinea:   zoneRequest.ColorLinea,
		ColorRelleno: zoneRequest.ColorRelleno,
		Puntos:       zoneRequest.Puntos,
		Nombre:       zoneRequest.Nombre,
	}
	tx := db.Save(&zone)

	if tx.Error != nil {
		return tx.Error
	}

	tx = db.Where("zona_id = ?", id).Delete(&model.ZonaVehiculo{})

	if len(zoneRequest.Imeis) > 0 {
		zonesWithVehicles := []model.ZonaVehiculo{}
		for _, imei := range zoneRequest.Imeis {
			zonesWithVehicles = append(zonesWithVehicles, model.ZonaVehiculo{
				ZonaID:        zone.ID,
				Imei:          imei,
				AvisarEntrada: zoneRequest.AvisarEntrada,
				AvisarSalida:  zoneRequest.AvisarSalida,
			})
		}

		tx = db.Create(&zonesWithVehicles)
	}

	return tx.Error
}

func (ma *zonasManager) DeleteZoneById(idParam string) error {
	db, close, err := db.ObtenerConexionDb()
	defer close()

	if err != nil {
		return err
	}

	id, idParseErr := strconv.Atoi(idParam)

	if idParseErr != nil {
		return idParseErr
	}

	zone := model.Zona{ID: int32(id)}
	tx := db.Delete(&zone)

	return tx.Error
}

func (ma *zonasManager) GetZoneConfigByImei(imei string) ([]model.ZoneView, error) {
	db, close, err := db.ObtenerConexionDb()
	defer close()

	zoneConfig := []model.ZoneView{}

	tx := db.Model(&model.ZonaVehiculo{}).Select("zona.puntos, zona.nombre, zona.id, zona_vehiculos.avisar_entrada, zona_vehiculos.avisar_salida").Joins("join zona on zona.id = zona_vehiculos.zona_id").Where("imei = ?", imei).Scan(&zoneConfig)

	if err != nil {
		return []model.ZoneView{}, err
	}

	return zoneConfig, tx.Error
}
