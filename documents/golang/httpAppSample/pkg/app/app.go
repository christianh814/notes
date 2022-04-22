package app

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

var ServerPort string = "8080"

func Start() {
	//create/register a new request multiplexer
	mux := http.NewServeMux()

	// define routes to that multiplexer
	mux.HandleFunc("/greet", greet)
	mux.HandleFunc("/customers", getAllCust)

	// log out that we're starting
	log.Info("Starting server on port " + ServerPort)

	// run and listen to on 8080 on that multiplexer, display any errors if it errors out
	log.Fatal(http.ListenAndServe(":"+ServerPort, mux))

}
