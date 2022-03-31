package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

type Customer struct {
	Name    string `json:"full_name" xml:"name"`
	City    string `json:"city" xml:"city"`
	Zipcode string `json:"zip_code" xml:"zip"`
}

func main() {
	// define route
	http.HandleFunc("/greet", greet)
	http.HandleFunc("/customers", getAllCust)

	// run and listen to on 8080
	_ = http.ListenAndServe(":8080", nil)
}

// greet prints hello world
func greet(w http.ResponseWriter, r *http.Request) {
	// this print hello world when someone vists the endpoint
	fmt.Fprint(w, "Hello World")
}

// getAllCust prints customer data as JSON
func getAllCust(w http.ResponseWriter, r *http.Request) {
	// Set dummy values
	cust := []Customer{
		{Name: "Christian Hernandez", City: "Los Angeles, CA", Zipcode: "90293"},
		{Name: "Kacie", City: "Oneill", Zipcode: "12345"},
	}

	//if someone requests XML, send it to them. If they request JSON, send that
	if r.Header.Get("Content-Type") == "application/xml" {

		// set the right header for XML
		w.Header().Add("Content-Type", "application/xml")

		// encode struct as XML and print it out
		xml.NewEncoder(w).Encode(cust)

	} else {
		// set the right header for JSON
		w.Header().Add("Content-Type", "application/json")

		// encode struct as JSON and print it out
		json.NewEncoder(w).Encode(cust)

	}

}
