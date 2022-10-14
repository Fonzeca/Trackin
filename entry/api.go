package entry

import (
	"github.com/Fonzeca/Trackin/entry/manager"
)

var (
	DataEntryManager = manager.NewDataEntryManager()
)

// func Router(e *echo.Echo) {
// 	e.POST("/data", dataEntryApi)
// }

// func dataEntryApi(c echo.Context) error {
// 	//Cacheo de raw body
// 	json_map := make(map[string]interface{})
// 	json.NewDecoder(c.Request().Body).Decode(&json_map)

// 	bytesJson, err := json.Marshal(json_map)
// 	if err != nil {
// 		return err
// 	}

// 	//Bind de datos que nos interesan
// 	data := jsonModel.SimplyData{}
// 	json.Unmarshal(bytesJson, &data)
// 	if err != nil {
// 		return err
// 	}

// 	//Seteamos el PayLoad
// 	data.PayLoad = string(bytesJson)

// 	//Mandamos los datos a guardar
// 	DataEntryManager.CanalEntrada <- data
// 	return c.NoContent(http.StatusOK)
// }
