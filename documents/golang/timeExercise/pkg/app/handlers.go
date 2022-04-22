package app

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	gorillamux "github.com/gorilla/mux"
)

type CurrentTime struct {
	TimeZone    string `json:"time_zone"`
	CurrentTime string `json:"current_time"`
}

func getTime(w http.ResponseWriter, r *http.Request) {
	// set response as JSON
	w.Header().Add("Content-Type", "application/json")

	// get the vars
	vars := gorillamux.Vars(r)

	// set up an empty currtime
	var currtime = []CurrentTime{}

	// Split location
	nl := strings.Split(vars["tz"], ",")
	for _, l := range nl {
		tz, err := time.LoadLocation(l)
		if err != nil {
			// Write 404 header
			w.WriteHeader(http.StatusNotFound)
			currtime = append(currtime, CurrentTime{TimeZone: "UTC", CurrentTime: time.Now().UTC().String()})
		} else {
			currtime = append(currtime, CurrentTime{TimeZone: l, CurrentTime: time.Now().In(tz).String()})

		}
	}

	json.NewEncoder(w).Encode(currtime)
}
