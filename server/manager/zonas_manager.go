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

// GetZonesByEmpresaId obtiene todas las zonas de una empresa, incluyendo las asociadas y no asociadas a vehículos.
// Retorna una lista de ZoneRequest, agrupando los imeis de los vehículos asociados a cada zona.
func (ma *zonasManager) GetZonesByEmpresaId(idParam string) ([]model.ZoneRequest, error) {
	db, close, err := db.ObtenerConexionDb()

	if err != nil {
		return nil, err
	}
	defer close()

	id, idParseErr := strconv.Atoi(idParam)
	if idParseErr != nil {
		return nil, idParseErr
	}

	// Consulta única con LEFT JOIN para traer zonas con y sin vehículos asociados
	zones := []model.ZoneView{}
	tx := db.Model(&model.Zona{}).
		Select("zona.id, zona.empresa_id, zona.color_linea, zona.color_relleno, zona.puntos, zona.nombre, zona.velocidad_maxima, zona_vehiculos.imei, zona_vehiculos.avisar_entrada, zona_vehiculos.avisar_salida").
		Joins("LEFT JOIN zona_vehiculos ON zona.id = zona_vehiculos.zona_id").
		Where("zona.empresa_id = ?", id).
		Order("zona.id desc").
		Scan(&zones)

	if tx.Error != nil {
		return nil, tx.Error
	}

	// Agrupar los resultados por zona
	zoneMap := make(map[int32]*model.ZoneRequest)
	for _, zone := range zones {
		zr, exists := zoneMap[zone.Id]
		if !exists {
			zr = &model.ZoneRequest{
				Id:              zone.Id,
				EmpresaId:       zone.EmpresaId,
				ColorLinea:      zone.ColorLinea,
				ColorRelleno:    zone.ColorRelleno,
				Puntos:          zone.Puntos,
				Nombre:          zone.Nombre,
				Imeis:           []string{},
				AvisarEntrada:   zone.AvisarEntrada,
				AvisarSalida:    zone.AvisarSalida,
				VelocidadMaxima: zone.VelocidadMaxima,
			}
			zoneMap[zone.Id] = zr
		}
		if zone.Imei != "" {
			zr.Imeis = append(zr.Imeis, zone.Imei)
			zr.AvisarEntrada = zone.AvisarEntrada
			zr.AvisarSalida = zone.AvisarSalida
		}
	}

	// Convertir el mapa a slice
	result := make([]model.ZoneRequest, 0, len(zoneMap))
	for _, v := range zoneMap {
		result = append(result, *v)
	}

	return result, nil
}

func (ma *zonasManager) CreateZone(zoneRequest model.ZoneRequest) error {
	db, close, err := db.ObtenerConexionDb()
	if err != nil {
		return err
	}
	defer close()

	transactionErr := db.Transaction(func(tx *gorm.DB) error {

		zone := model.Zona{
			EmpresaID:       int32(zoneRequest.EmpresaId),
			ColorLinea:      zoneRequest.ColorLinea,
			ColorRelleno:    zoneRequest.ColorRelleno,
			Puntos:          zoneRequest.Puntos,
			Nombre:          zoneRequest.Nombre,
			VelocidadMaxima: zoneRequest.VelocidadMaxima,
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
	if err != nil {
		return err
	}
	defer close()

	id, idParseErr := strconv.Atoi(idParam)

	if idParseErr != nil {
		return idParseErr
	}

	zone := model.Zona{
		ID:              int32(id),
		EmpresaID:       int32(zoneRequest.EmpresaId),
		ColorLinea:      zoneRequest.ColorLinea,
		ColorRelleno:    zoneRequest.ColorRelleno,
		Puntos:          zoneRequest.Puntos,
		Nombre:          zoneRequest.Nombre,
		VelocidadMaxima: zoneRequest.VelocidadMaxima,
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
	if err != nil {
		return err
	}
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
	if err != nil {
		return []model.ZoneView{}, err
	}
	defer close()

	zoneConfig := []model.ZoneView{}

	tx := db.Model(&model.ZonaVehiculo{}).
		Select("zona.puntos, zona.nombre, zona.id, zona_vehiculos.avisar_entrada, zona_vehiculos.avisar_salida, zona.velocidad_maxima").
		Joins("join zona on zona.id = zona_vehiculos.zona_id").
		Where("imei = ?", imei).
		Scan(&zoneConfig)

	if err != nil {
		return []model.ZoneView{}, err
	}

	return zoneConfig, tx.Error
}
