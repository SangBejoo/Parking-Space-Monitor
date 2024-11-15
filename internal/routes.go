// internal/routes.go
package main

import (
	"net/http"

	"github.com/SangBejoo/parking-space-monitor/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	taxiHandler := &handlers.TaxiHandler{}
	placeHandler := &handlers.PlaceHandler{}

	// Define the bulk insert route before the dynamic {id} routes
	router.HandleFunc("/taxis/bulk", taxiHandler.CreateMultipleTaxiLocations).Methods("POST")

	router.HandleFunc("/taxis", taxiHandler.CreateTaxiLocation).Methods("POST")
	router.HandleFunc("/taxis", taxiHandler.GetAllTaxiLocations).Methods("GET")
	router.HandleFunc("/taxis/{id}", taxiHandler.GetTaxiLocation).Methods("GET")
	router.HandleFunc("/taxis/{id}", taxiHandler.UpdateTaxiLocation).Methods("PUT")
	router.HandleFunc("/taxis/{id}", taxiHandler.DeleteTaxiLocation).Methods("DELETE")

	router.HandleFunc("/places", placeHandler.CreatePlace).Methods("POST")
	router.HandleFunc("/places", placeHandler.GetAllPlaces).Methods("GET")
	router.HandleFunc("/places/{id}", placeHandler.GetPlace).Methods("GET")
	router.HandleFunc("/places/{id}", placeHandler.UpdatePlace).Methods("PUT")
	router.HandleFunc("/places/{id}", placeHandler.DeletePlace).Methods("DELETE")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
