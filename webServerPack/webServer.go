package webserverpack

import (
	"encoding/json"
	"fmt"
	"net/http"
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

/*
  Questo webServer espone l'endpoint /metrics, il quale effettua una GET verso l'indirizzo del ESP-Moisture-Monitor/allResources
*/

func AskMetricToESP(w http.ResponseWriter, r *http.Request) {
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

	fmt.Fprintf(w, "OUTPUT: %s\n", data.String())
}

func DefineHTTPWebServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ciao")
	})

	http.HandleFunc("/metrics", AskMetricToESP)

	fmt.Println("Server in ascolto su http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Errore nell'avvio del server:", err)
	}
}
