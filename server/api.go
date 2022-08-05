package server

import (
	"net/http"

	"github.com/Fonzeca/Trackin/server/manager"
	"github.com/labstack/echo/v4"
)

func Router(c *echo.Echo) error {
	return nil
}

type api struct {
	manager manager.Manager
}

func NewApi() api {
	m := manager.NewManager()
	return api{manager: m}
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
