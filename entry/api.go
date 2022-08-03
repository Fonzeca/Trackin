package entry

import (
	"net/http"

	jsonModel "github.com/Fonzeca/Trackin/entry/json"
	"github.com/Fonzeca/Trackin/entry/manager"
	"github.com/labstack/echo/v4"
)

var (
	DataEntryManager = manager.NewDataEntryManager()
)

func Router(e *echo.Echo) {
	e.POST("/data", dataEntryApi)
}

func dataEntryApi(c echo.Context) error {
	data := jsonModel.SimplyData{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}
	DataEntryManager.CanalEntrada <- data
	return c.NoContent(http.StatusOK)
}
