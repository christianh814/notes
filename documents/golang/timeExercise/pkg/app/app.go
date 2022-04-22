package app

import (
	"net/http"

	gorillamux "github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var ServerPort string = "8080"

func Start() {
	//define mux router
	router := gorillamux.NewRouter()

	// define GET routes
	router.HandleFunc("/api/time", getTime).Queries("tz", "{tz:[0-9A-Za-z_/]+}").Methods(http.MethodGet)
	router.HandleFunc("/api/time", getTime).Methods(http.MethodGet)

	// log out that we are starting
	log.Info("Starting on port " + ServerPort)

	// run/listen on port and log if there's an error
	log.Fatal(http.ListenAndServe(":"+ServerPort, router))
}
