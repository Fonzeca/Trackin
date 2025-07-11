package suites

import (
	"time"

	"github.com/Fonzeca/Trackin/internal/infrastructure/database/model"
	modeljson "github.com/Fonzeca/Trackin/internal/interfaces/messaging/json"
)

type SuiteZonaProccessData struct {
	Nombre                 string
	Zonas                  []model.ZoneView
	Delay                  time.Duration
	FinishDelay            time.Duration
	Data                   []modeljson.SimplyData
	ExpectedEventTypeCalls []model.ZoneNotification
}

var ZoneProcces_NormalSuite = SuiteZonaProccessData{
	Nombre: "NormalDemo",
	Zonas: []model.ZoneView{
		{Puntos: "-48,-67; -49,-67; -49,-66; -48,-66", Id: 22, Nombre: "NombreZona - NormalDemo", AvisarEntrada: true, AvisarSalida: true},
	},
	Delay:       time.Millisecond * 200,
	FinishDelay: time.Second * 3,
	Data: func() (res []modeljson.SimplyData) {
		now := time.Now()
		data := modeljson.SimplyData{
			EngineStatus: true,
			Imei:         "3322",
			Latitude:     -48.5,
			Longitude:    -66.5,
			Date:         now,
		}
		res = append(res, data)
		data = modeljson.SimplyData{
			EngineStatus: true,
			Imei:         "4444",
			Latitude:     -48.5,
			Longitude:    -66.5,
			Date:         now,
		}
		res = append(res, data)
		data = modeljson.SimplyData{
			EngineStatus: true,
			Imei:         "3322",
			Latitude:     55,
			Longitude:    55,
			Date:         now.Add(5 * time.Second),
		}
		res = append(res, data)
		data = modeljson.SimplyData{
			EngineStatus: true,
			Imei:         "4444",
			Latitude:     55,
			Longitude:    55,
			Date:         now.Add(5 * time.Second),
		}
		res = append(res, data)
		data = modeljson.SimplyData{
			EngineStatus: true,
			Imei:         "3322",
			Latitude:     -48.5,
			Longitude:    -66.5,
			Date:         now.Add(10 * time.Second),
		}
		res = append(res, data)
		data = modeljson.SimplyData{
			EngineStatus: true,
			Imei:         "4444",
			Latitude:     -48.5,
			Longitude:    -66.5,
			Date:         now.Add(10 * time.Second),
		}
		res = append(res, data)
		return res
	}(),
	ExpectedEventTypeCalls: []model.ZoneNotification{
		{Imei: "3322", ZoneName: "NombreZona - NormalDemo", ZoneID: 22, EventType: "sale"},
		{Imei: "4444", ZoneName: "NombreZona - NormalDemo", ZoneID: 22, EventType: "sale"},
		{Imei: "3322", ZoneName: "NombreZona - NormalDemo", ZoneID: 22, EventType: "entra"},
		{Imei: "4444", ZoneName: "NombreZona - NormalDemo", ZoneID: 22, EventType: "entra"},
	},
}

var ZoneProcces_EstressSuite = SuiteZonaProccessData{
	Nombre: "EstressDemo",
	Zonas: []model.ZoneView{
		{Puntos: "-48,-67; -49,-67; -49,-66; -48,-66", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.285119,-67.811417; -49.263839,-67.779123; -49.277146,-67.757615; -49.235108,-67.759402; -49.231656,-67.816017", Id: 22, Nombre: "ZonaSalidaSanju - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
		{Puntos: "-49.481866,-68.053486; -49.156023,-67.858349; -49.292352,-67.418602; -48.900296,-67.597249; -48.923763,-68.144184", Id: 22, Nombre: "NombreZona - EstressDemo", AvisarEntrada: true, AvisarSalida: true},
	},
	Delay:       time.Millisecond * 500,
	FinishDelay: time.Second * 3,
	Data: func() (res []modeljson.SimplyData) {
		now := time.Now()
		data := modeljson.SimplyData{
			EngineStatus: true,
			Imei:         "3322",
			Latitude:     -49.293650,
			Longitude:    -67.779043,
			Date:         now,
		}
		res = append(res, data)
		data = modeljson.SimplyData{
			EngineStatus: true,
			Imei:         "3322",
			Latitude:     -49.283641,
			Longitude:    -67.778567,
			Date:         now.Add(5 * time.Second),
		}
		res = append(res, data)
		data = modeljson.SimplyData{
			EngineStatus: true,
			Imei:         "3322",
			Latitude:     -49.256535,
			Longitude:    -67.778224,
			Date:         now.Add(10 * time.Second),
		}
		res = append(res, data)
		data = modeljson.SimplyData{
			EngineStatus: true,
			Imei:         "3322",
			Latitude:     -49.291478,
			Longitude:    -67.782346,
			Date:         now.Add(15 * time.Second),
		}
		res = append(res, data)
		data = modeljson.SimplyData{
			EngineStatus: true,
			Imei:         "3322",
			Latitude:     -49.256983,
			Longitude:    -67.779254,
			Date:         now.Add(20 * time.Second),
		}
		res = append(res, data)
		return res
	}(),
	ExpectedEventTypeCalls: []model.ZoneNotification{
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - EstressDemo", ZoneID: 22, EventType: "entra"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - EstressDemo", ZoneID: 22, EventType: "sale"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - EstressDemo", ZoneID: 22, EventType: "entra"},
	},
}

var ZoneProcces_PrecisionSuite = SuiteZonaProccessData{
	Nombre: "PrecisionDemo",
	Zonas: []model.ZoneView{
		{Puntos: "-49.285119,-67.811417; -49.263839,-67.779123; -49.277146,-67.757615; -49.235108,-67.759402; -49.231656,-67.816017", Id: 22, Nombre: "ZonaSalidaSanju - PrecisionDemo", AvisarEntrada: true, AvisarSalida: true},
	},
	Delay:       time.Millisecond * 200,
	FinishDelay: time.Second * 3,
	Data: func() (res []modeljson.SimplyData) {
		now := time.Now()
		puntos := []struct {
			Lat float64
			Lng float64
		}{
			{Lng: -67.7995491, Lat: -49.2822801},
			{Lng: -67.8024673, Lat: -49.2748849},
			{Lng: -67.7940559, Lat: -49.2776862},
			{Lng: -67.7949142, Lat: -49.2704024},
			{Lng: -67.7870178, Lat: -49.2742125},
			{Lng: -67.7883911, Lat: -49.2679369},
			{Lng: -67.7825546, Lat: -49.2707386},
			{Lng: -67.7799797, Lat: -49.2668162},
			{Lng: -67.7741432, Lat: -49.2678248},
			{Lng: -67.7708817, Lat: -49.2700662},
			{Lng: -67.7629852, Lat: -49.2753331},
			{Lng: -67.7586079, Lat: -49.275713},
			{Lng: -67.7570629, Lat: -49.2742002},
			{Lng: -67.7589512, Lat: -49.2708384},
			{Lng: -67.7453899, Lat: -49.2625448},
			{Lng: -67.7543678, Lat: -49.2526455},
			{Lng: -67.7420082, Lat: -49.249058},
			{Lng: -67.7499046, Lat: -49.2463671},
			{Lng: -67.7447548, Lat: -49.2413214},
			{Lng: -67.7557411, Lat: -49.2423306},
			{Lng: -67.7643242, Lat: -49.2402001},
			{Lng: -67.7674141, Lat: -49.2454702},
			{Lng: -67.7763405, Lat: -49.2459187},
			{Lng: -67.7777138, Lat: -49.240985},
			{Lng: -67.7722206, Lat: -49.237733},
			{Lng: -67.7639122, Lat: -49.2485422},
			{Lng: -67.7549858, Lat: -49.2390114},
			{Lng: -67.7573891, Lat: -49.2336285},
			{Lng: -67.7826233, Lat: -49.2329556},
			{Lng: -67.7848549, Lat: -49.2327313},
			{Lng: -67.7881165, Lat: -49.2322827},
			{Lng: -67.792408, Lat: -49.2321705},
			{Lng: -67.8102608, Lat: -49.2277964},
			{Lng: -67.8104324, Lat: -49.2381143},
			{Lng: -67.8209038, Lat: -49.2360958},
			{Lng: -67.8164406, Lat: -49.2421512},
			{Lng: -67.8164406, Lat: -49.2523541},
			{Lng: -67.8150673, Lat: -49.2664776},
			{Lng: -67.8260536, Lat: -49.270288},
			{Lng: -67.8109474, Lat: -49.2717449},
			{Lng: -67.8260536, Lat: -49.2784683},
			{Lng: -67.8068275, Lat: -49.275667},
			{Lng: -67.823307, Lat: -49.2840704},
		}
		for _, point := range puntos {
			now = now.Add(5 * time.Second)
			data := modeljson.SimplyData{
				EngineStatus: true,
				Imei:         "3322",
				Latitude:     point.Lat,
				Longitude:    point.Lng,
				Date:         now,
			}
			res = append(res, data)
		}
		return res
	}(),
	ExpectedEventTypeCalls: []model.ZoneNotification{
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - PrecisionDemo", ZoneID: 22, EventType: "entra"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - PrecisionDemo", ZoneID: 22, EventType: "sale"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - PrecisionDemo", ZoneID: 22, EventType: "entra"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - PrecisionDemo", ZoneID: 22, EventType: "sale"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - PrecisionDemo", ZoneID: 22, EventType: "entra"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - PrecisionDemo", ZoneID: 22, EventType: "sale"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - PrecisionDemo", ZoneID: 22, EventType: "entra"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - PrecisionDemo", ZoneID: 22, EventType: "sale"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - PrecisionDemo", ZoneID: 22, EventType: "entra"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - PrecisionDemo", ZoneID: 22, EventType: "sale"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - PrecisionDemo", ZoneID: 22, EventType: "entra"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - PrecisionDemo", ZoneID: 22, EventType: "sale"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - PrecisionDemo", ZoneID: 22, EventType: "entra"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - PrecisionDemo", ZoneID: 22, EventType: "sale"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - PrecisionDemo", ZoneID: 22, EventType: "entra"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - PrecisionDemo", ZoneID: 22, EventType: "sale"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - PrecisionDemo", ZoneID: 22, EventType: "entra"},
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - PrecisionDemo", ZoneID: 22, EventType: "sale"},
	},
}

var ZoneProcces_EngineStatusSuite = SuiteZonaProccessData{
	Nombre: "EngineStatuDemo",
	Zonas: []model.ZoneView{
		{Puntos: "-49.285119,-67.811417; -49.263839,-67.779123; -49.277146,-67.757615; -49.235108,-67.759402; -49.231656,-67.816017", Id: 22, Nombre: "ZonaSalidaSanju - EngineStatusDemo", AvisarEntrada: true, AvisarSalida: true},
	},
	Delay:       time.Millisecond * 200,
	FinishDelay: time.Second * 3,
	Data: func() (res []modeljson.SimplyData) {
		now := time.Now()
		data := modeljson.SimplyData{
			EngineStatus: true,
			Imei:         "3322",
			Latitude:     -49.2791501,
			Longitude:    -67.7849579,
			Date:         now,
		}
		res = append(res, data)
		data = modeljson.SimplyData{
			EngineStatus: false,
			Imei:         "3322",
			Latitude:     -49.2561751,
			Longitude:    -67.7918243,
			Date:         now.Add(5 * time.Second),
		}
		res = append(res, data)
		data = modeljson.SimplyData{
			EngineStatus: true,
			Imei:         "3322",
			Latitude:     -49.2791501,
			Longitude:    -67.7849579,
			Date:         now.Add(10 * time.Second),
		}
		res = append(res, data)

		return res
	}(),
	ExpectedEventTypeCalls: []model.ZoneNotification{
		{Imei: "3322", ZoneName: "ZonaSalidaSanju - EngineStatusDemo", ZoneID: 22, EventType: "entra"},
	},
}
