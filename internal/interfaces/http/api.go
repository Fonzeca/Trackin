package server

import (
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/Fonzeca/Trackin/internal/core"
	manager "github.com/Fonzeca/Trackin/internal/core/managers"
	"github.com/Fonzeca/Trackin/internal/infrastructure/database/model"
	"github.com/labstack/echo/v4"
)

func Router(c *echo.Echo) error {
	return nil
}

type api struct {
	container *manager.ManagerContainer
}

func NewApi() api {
	container := manager.GetManagerContainer()
	return api{container: container}
}

func (api *api) GetLastLogByImei(c echo.Context) error {

	val, paramErr := c.FormParams()
	imei := val.Get("imei")

	//Si esto va atener un frontend inteligente, borramos todos los if de errores
	if paramErr != nil {
		return c.JSON(http.StatusBadRequest, paramErr.Error())
	}

	if imei == "" {
		return c.JSON(http.StatusBadRequest, "Parámetro imei incorrecto")
	}

	log, logErr := api.container.GetRoutesManager().GetLastLogByImei(imei)

	if logErr != nil {
		return c.JSON(http.StatusNotFound, logErr.Error())
	}

	return c.JSON(http.StatusOK, log)
}

func (api *api) GetVehiclesStateByImeis(c echo.Context) error {
	data := model.ImeisBody{}
	c.Bind(&data)

	logs, logErr := api.container.GetRoutesManager().GetVehiclesStateByImeis(data)

	if logErr != nil {
		return c.JSON(http.StatusNotFound, logErr.Error())
	}

	return c.JSON(http.StatusOK, logs)
}

func (api *api) GetRouteByImei(c echo.Context) error {
	data := model.RouteRequest{}
	c.Bind(&data)

	if len(data.ZonesIds) > 0 {
		zones, zoneErr := api.container.GetZonasManager().GetZoneByIds(data.ZonesIds)
		if zoneErr != nil {
			return c.JSON(http.StatusNotFound, zoneErr.Error())
		}

		// Si se encontraron zonas, las pasamos a la función de la ruta
		route, logErr := api.container.GetRoutesManager().GetRouteByImeiAndZones(data, zones)
		if logErr != nil {
			return c.JSON(http.StatusNotFound, logErr.Error())
		}

		return c.JSON(http.StatusOK, route)
	} else {
		// Si no se especifican zonas, obtenemos la ruta por IMEI sin zonas
		route, logErr := api.container.GetRoutesManager().GetRouteByImei(data)

		if logErr != nil {
			return c.JSON(http.StatusNotFound, logErr.Error())
		}

		return c.JSON(http.StatusOK, route)
	}
}

func (api *api) GetZonesByEmpresaId(c echo.Context) error {
	val, _ := c.FormParams()
	id := val.Get("id")

	zones, zoneErr := api.container.GetZonasManager().GetZonesWithImeisByEmpresaId(id)

	if zoneErr != nil {
		return c.JSON(http.StatusNotFound, zoneErr.Error())
	}

	return c.JSON(http.StatusOK, zones)
}

func (api *api) CreateZone(c echo.Context) error {
	data := model.ZoneRequest{}
	c.Bind(&data)

	zoneErr := api.container.GetZonasManager().CreateZone(data)

	if zoneErr != nil {
		return c.JSON(http.StatusBadRequest, zoneErr.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (api *api) EditZoneById(c echo.Context) error {
	data := model.ZoneRequest{}
	c.Bind(&data)

	val, _ := c.FormParams()
	id := val.Get("id")

	zoneErr := api.container.GetZonasManager().EditZoneById(id, data)

	if zoneErr != nil {
		return c.JSON(http.StatusBadRequest, zoneErr.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (api *api) DeleteZoneById(c echo.Context) error {
	val, _ := c.FormParams()
	id := val.Get("id")

	zoneErr := api.container.GetZonasManager().DeleteZoneById(id)

	if zoneErr != nil {
		return c.JSON(http.StatusBadRequest, zoneErr.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (api *api) GetSummaryRoutesAndZones(c echo.Context) error {
	data := model.SummaryRequest{}
	c.Bind(&data)

	intersections, err := api.container.GetRoutesManager().GetSummaryRoutesAndZones(&data)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	records, err := core.MapIntersectionsToCSV(intersections)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error generating CSV: "+err.Error())
	}

	filename := fmt.Sprintf("summary_%s.csv", data.Imei)

	// Set headers for CSV download
	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	// Create a new CSV writer
	writer := csv.NewWriter(c.Response().Writer)
	defer writer.Flush() // Ensure all data is written

	// Write data to the CSV
	if err := writer.WriteAll(records); err != nil {
		return c.String(http.StatusInternalServerError, "Error writing CSV: "+err.Error())
	}

	return nil // No explicit response body needed as data is written directly to writer
}
