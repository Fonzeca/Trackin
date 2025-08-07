package manager

import (
	"strconv"

	db "github.com/Fonzeca/Trackin/internal/infrastructure/database"
	"github.com/Fonzeca/Trackin/internal/infrastructure/database/model"
	"gorm.io/gorm"
)

type IZonasManager interface {
	GetZonesWithImeisByEmpresaId(idParam string) ([]model.ZoneRequest, error)
	GetZonesByEmpresaId(empresaId int32) ([]model.Zona, error)
	CreateZone(zoneRequest model.ZoneRequest) error
	EditZoneById(idParam string, zoneRequest model.ZoneRequest) error
	DeleteZoneById(idParam string) error
	GetZoneConfigByImei(imei string) ([]model.ZoneView, error)
	GetZoneByIds(ids []int32) ([]model.ZoneView, error)

	// Setter para inyección de dependencias
	SetRoutesManager(routesManager IRoutesManager)
}

type zonasManager struct {
	routesManager IRoutesManager
}

func newZonasManager() IZonasManager {
	return &zonasManager{}
}

// SetRoutesManager inyecta la dependencia del routes manager
func (ma *zonasManager) SetRoutesManager(routesManager IRoutesManager) {
	ma.routesManager = routesManager
}

func (ma *zonasManager) GetZonesByEmpresaId(empresaId int32) ([]model.Zona, error) {
	zones := []model.Zona{}
	tx := db.DB.Model(&model.Zona{}).Where("empresa_id = ?", empresaId).Find(&zones)

	if tx.Error != nil {
		return nil, tx.Error
	}

	if len(zones) == 0 {
		return nil, nil // No zones found
	}

	return zones, nil
}

// GetZonesByEmpresaId obtiene todas las zonas de una empresa, incluyendo las asociadas y no asociadas a vehículos.
// Retorna una lista de ZoneRequest, agrupando los imeis de los vehículos asociados a cada zona.
func (ma *zonasManager) GetZonesWithImeisByEmpresaId(idParam string) ([]model.ZoneRequest, error) {
	id, idParseErr := strconv.Atoi(idParam)
	if idParseErr != nil {
		return nil, idParseErr
	}

	// Consulta única con LEFT JOIN para traer zonas con y sin vehículos asociados
	zones := []model.ZoneView{}
	tx := db.DB.Model(&model.Zona{}).
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
	zoneIds := make([]int32, 0, len(zones))
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
			zoneIds = append(zoneIds, zone.Id)
		}
		if zone.Imei != "" {
			zr.Imeis = append(zr.Imeis, zone.Imei)
			zr.AvisarEntrada = zone.AvisarEntrada
			zr.AvisarSalida = zone.AvisarSalida
		}
	}

	// Convertir el mapa a slice
	result := make([]model.ZoneRequest, 0, len(zoneMap))
	for _, v := range zoneIds {
		result = append(result, *zoneMap[v])
	}

	return result, nil
}

func (ma *zonasManager) CreateZone(zoneRequest model.ZoneRequest) error {
	transactionErr := db.DB.Transaction(func(tx *gorm.DB) error {

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
	tx := db.DB.Save(&zone)

	if tx.Error != nil {
		return tx.Error
	}

	tx = db.DB.Where("zona_id = ?", id).Delete(&model.ZonaVehiculo{})

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

		tx = db.DB.Create(&zonesWithVehicles)
	}

	return tx.Error
}

func (ma *zonasManager) GetZoneByIds(ids []int32) ([]model.ZoneView, error) {
	var zones []model.ZoneView
	tx := db.DB.Where("id IN ?", ids).Find(&zones)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return zones, nil
}

func (ma *zonasManager) DeleteZoneById(idParam string) error {
	id, idParseErr := strconv.Atoi(idParam)

	if idParseErr != nil {
		return idParseErr
	}

	zone := model.Zona{ID: int32(id)}
	tx := db.DB.Delete(&zone)

	return tx.Error
}

func (ma *zonasManager) GetZoneConfigByImei(imei string) ([]model.ZoneView, error) {
	zoneConfig := []model.ZoneView{}

	tx := db.DB.Model(&model.ZonaVehiculo{}).
		Select("zona.puntos, zona.nombre, zona.id, zona_vehiculos.avisar_entrada, zona_vehiculos.avisar_salida, zona.velocidad_maxima").
		Joins("join zona on zona.id = zona_vehiculos.zona_id").
		Where("imei = ?", imei).
		Scan(&zoneConfig)

	return zoneConfig, tx.Error
}
