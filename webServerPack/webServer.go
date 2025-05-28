package webserverpack

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const ESP_ADDRESS = "192.168.1.169"

// [{"i":0,"v":3228,"n":41,"d":5392,"w":2808,"t":4651}
type moistureSensor struct {
	I int `json:"i"`
	V int `json:"v"`
	N int `json:"n"`
	D int `json:"d"`
	W int `json:"w"`
	T int `json:"t"`
}

func (m moistureSensor) String() string {
	return fmt.Sprint("sensor: ", m.I, " value: ", m.V, " plant: ", m.N, " dry: ", m.D, " wet: ", m.W, " treshold: ", m.T)
}

type pumpState struct {
	I int  `json:"i"`
	S bool `json:"s"`
}

func (p pumpState) String() string {
	var state string
	if p.S {
		state = "ON"
	} else {
		state = "OFF"
	}
	return fmt.Sprint("pump: ", p.I, " state: ", state)
}

type AllResources struct {
	Ps []pumpState      `json:"pS"`
	Ms []moistureSensor `json:"mS"`
}

func (aR AllResources) String() string {
	var msg string
	for _, ps := range aR.Ps {
		msg += ps.String() + "\n"
	}
	msg += "---------\n"
	for _, ms := range aR.Ms {
		msg += ms.String() + "\n"
	}
	return msg
}

var (
	pumpStateProm = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pump_state",
			Help: "Pump state (1 = ON, 0 = OFF)",
		},
		[]string{"pump"},
	)

	sensorValue = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sensor_value",
			Help: "Current moisture sensor value",
		},
		[]string{"sensor", "plant", "connected"},
	)
	sensorDry = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sensor_dry",
			Help: "Dry calibration value",
		},
		[]string{"sensor", "plant", "connected"},
	)

	sensorWet = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sensor_wet",
			Help: "Wet calibration value",
		},
		[]string{"sensor", "plant", "connected"},
	)

	sensorThreshold = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sensor_threshold",
			Help: "Threshold value",
		},
		[]string{"sensor", "plant", "connected"},
	)
)

func init() {
	prometheus.MustRegister(pumpStateProm)
	prometheus.MustRegister(sensorValue)
	prometheus.MustRegister(sensorDry)
	prometheus.MustRegister(sensorWet)
	prometheus.MustRegister(sensorThreshold)
}

/*
  Questo webServer espone l'endpoint /metrics, il quale effettua una GET verso l'indirizzo del ESP-Moisture-Monitor/allResources
*/

func AskMetricToESP(w io.Writer, r *http.Request) {
	resp, err := http.Get("http://" + ESP_ADDRESS + "/allResources") // endpoint JSON
	if err != nil {
		fmt.Fprintf(w, "Errore durante la GET -> %s\n", err.Error())
		return
	}
	defer resp.Body.Close()

	var data AllResources
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Fprintf(w, "Errore nel parsing del JSON: %s\n", err.Error())
		return
	}

	for _, pump := range data.Ps {
		val := 0.0
		if pump.S {
			val = 1.0
		}
		pumpStateProm.With(prometheus.Labels{"pump": fmt.Sprintf("%d", pump.I)}).Set(val)
	}

	for _, sensor := range data.Ms {
		connected := "true"
		if sensor.V < 0 {
			connected = "false"
		}

		labels := prometheus.Labels{
			"sensor":    fmt.Sprintf("%d", sensor.I),
			"connected": connected,
			"plant":     fmt.Sprintf("%d", sensor.N),
		}

		sensorValue.With(labels).Set(float64(sensor.V))
		sensorDry.With(labels).Set(float64(sensor.D))
		sensorWet.With(labels).Set(float64(sensor.W))
		sensorThreshold.With(labels).Set(float64(sensor.T))
	}

	fmt.Printf("OUTPUT: %s\n", data.String())
}

func DefineHTTPWebServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ciao")
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		AskMetricToESP(os.Stdout, nil)

		promhttp.Handler().ServeHTTP(w, r)
	})

	fmt.Println("Server in ascolto su http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Errore nell'avvio del server:", err)
	}
}
