package server

import (
	"net/http"

	"github.com/Fonzeca/Trackin/db/model"
	"github.com/Fonzeca/Trackin/server/manager"
	"github.com/labstack/echo/v4"
)

func Router(c *echo.Echo) error {
	return nil
}

type api struct {
	manager      manager.Manager
	zonasManager manager.IZonasManager
}

func NewApi() api {
	m := manager.NewManager()
	zonasManager := manager.ZonasManager
	return api{manager: m, zonasManager: zonasManager}
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

	log, logErr := api.manager.GetLastLogByImei(imei)

	if logErr != nil {
		return c.JSON(http.StatusNotFound, logErr.Error())
	}

	return c.JSON(http.StatusOK, log)
}

func (api *api) GetVehiclesStateByImeis(c echo.Context) error {
	data := model.ImeisBody{}
	c.Bind(&data)

	val, _ := c.FormParams()
	only := val.Get("only")

	logs, logErr := api.manager.GetVehiclesStateByImeis(only, data)

	if logErr != nil {
		return c.JSON(http.StatusNotFound, logErr.Error())
	}

	return c.JSON(http.StatusOK, logs)
}

func (api *api) GetRouteByImei(c echo.Context) error {
	data := model.RouteRequest{}
	c.Bind(&data)

	route, logErr := api.manager.GetRouteByImei(data)

	if logErr != nil {
		return c.JSON(http.StatusNotFound, logErr.Error())
	}

	return c.JSON(http.StatusOK, route)
}

func (api *api) GetZonesByEmpresaId(c echo.Context) error {
	val, _ := c.FormParams()
	id := val.Get("id")

	zones, zoneErr := api.zonasManager.GetZonesByEmpresaId(id)

	if zoneErr != nil {
		return c.JSON(http.StatusNotFound, zoneErr.Error())
	}

	return c.JSON(http.StatusOK, zones)
}

func (api *api) CreateZone(c echo.Context) error {
	data := model.ZoneRequest{}
	c.Bind(&data)

	zoneErr := api.zonasManager.CreateZone(data)

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

	zoneErr := api.zonasManager.EditZoneById(id, data)

	if zoneErr != nil {
		return c.JSON(http.StatusBadRequest, zoneErr.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (api *api) DeleteZoneById(c echo.Context) error {
	val, _ := c.FormParams()
	id := val.Get("id")

	zoneErr := api.zonasManager.DeleteZoneById(id)

	if zoneErr != nil {
		return c.JSON(http.StatusBadRequest, zoneErr.Error())
	}

	return c.NoContent(http.StatusOK)
}
