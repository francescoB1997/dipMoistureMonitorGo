package webserverpack

import (
	"fmt"
	"net/http"
)

func DefineHTTPWebServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ciao")
	})

	fmt.Println("Server in ascolto su http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Errore nell'avvio del server:", err)
	}
}
