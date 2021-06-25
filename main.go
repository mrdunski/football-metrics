package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
)

func main() {
	address := getListeningAddress()
	fmt.Printf("Listening on %s\n", address)

	http.HandleFunc("/", home)
	http.HandleFunc("/health", health)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(address, nil))
}

func getListeningAddress() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":8080"
	}

	return ":" + port
}

func home(w http.ResponseWriter, r *http.Request) {
	writeResponse("Hello!", w)
}

func health(w http.ResponseWriter, r *http.Request) {
	writeResponse("ok", w)
}

func writeResponse(r string, w http.ResponseWriter) {
	_, err := w.Write([]byte(r))
	if err != nil {
		log.Printf("Couldn't reply: %s\n", err)
	}
}
