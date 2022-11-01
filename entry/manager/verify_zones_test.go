package manager

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Fonzeca/Trackin/mocks"
	"github.com/Fonzeca/Trackin/server/manager"
	"github.com/Fonzeca/Trackin/services"
	"github.com/Fonzeca/Trackin/test/suites"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestZonaProccesData(t *testing.T) {
	testSuites := []suites.SuiteZonaProccessData{
		suites.ZoneProcces_NormalSuite,
		suites.ZoneProcces_EstressSuite,
		suites.ZoneProcces_PrecisionSuite,
		suites.ZoneProcces_EngineStatusSuite,
	}

	for _, suite := range testSuites {
		t.Run(suite.Nombre, func(t *testing.T) {
			zonasMockeado := mocks.NewIZonasManager(t)
			zonasMockeado.On("GetZoneConfigByImei", mock.Anything).Return(suite.Zonas, nil)
			manager.ZonasManager = zonasMockeado

			// ------------------------------------------------

			senderMoackeado := mocks.NewISender(t)
			senderMoackeado.On("SendMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			services.GlobalSender = senderMoackeado

			t.Log("Empezando test...")

			geoFenceDetector := NewGeofenceDetector()

			for _, data := range suite.Data {
				geoFenceDetector.DispatchMessage(data)
				time.Sleep(suite.Delay)
			}
			time.Sleep(suite.FinishDelay)

			assert.Equal(t, len(suite.ExpectedEventTypeCalls), len(senderMoackeado.Calls))

			for i, c := range senderMoackeado.Calls {
				expected := suite.ExpectedEventTypeCalls[i]
				zoneNotificationBytes, _ := json.Marshal(expected)

				c.Arguments.Assert(t, mock.Anything, "notification.zone.back.preparing", zoneNotificationBytes)
			}

			t.Log("Finish!")
		})
	}
}
