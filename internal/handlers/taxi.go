// internal/handlers/taxi.go
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"database/sql"

	"github.com/SangBejoo/parking-space-monitor/internal/models"
	"github.com/SangBejoo/parking-space-monitor/internal/repository"
	"github.com/gorilla/mux"
)

// TaxiHandler handles HTTP requests for Taxi operations.
type TaxiHandler struct {
	Repo *repository.TaxiRepository
}

// CreateTaxi handles the creation of a new taxi.
func (th *TaxiHandler) CreateTaxi(w http.ResponseWriter, r *http.Request) {
	var location models.TaxiLocation
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := th.Repo.CreateTaxi(location); err != nil {
		http.Error(w, "Failed to create taxi location", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"taxi_id":   location.TaxiID,
		"longitude": location.Longitude,
		"latitude":  location.Latitude,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetAllTaxis retrieves all taxis.
func (th *TaxiHandler) GetAllTaxis(w http.ResponseWriter, r *http.Request) {
	taxis, err := th.Repo.GetAllTaxis()
	if err != nil {
		http.Error(w, "Failed to query taxi locations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(taxis)
}

// GetTaxi retrieves a single taxi by ID.
func (th *TaxiHandler) GetTaxi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taxiID := vars["id"]

	taxi, err := th.Repo.GetTaxiByID(taxiID)
	if err == sql.ErrNoRows {
		http.Error(w, "Taxi not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to query taxi location", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(taxi)
}

// UpdateTaxi handles updating an existing taxi.
func (th *TaxiHandler) UpdateTaxi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taxiID := vars["id"]

	var location models.TaxiLocation
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := th.Repo.UpdateTaxi(taxiID, location); err != nil {
		if err.Error() == "taxi not found" {
			http.Error(w, "Taxi not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update taxi location", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Taxi location updated.")
}

// DeleteTaxi handles deleting a taxi by ID.
func (th *TaxiHandler) DeleteTaxi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taxiID := vars["id"]

	log.Printf("Received DELETE request for Taxi ID: %s", taxiID)

	if err := th.Repo.DeleteTaxi(taxiID); err != nil {
		log.Printf("Error deleting Taxi ID %s: %v", taxiID, err)
		if err.Error() == "taxi not found" {
			http.Error(w, "Taxi not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete taxi location", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully deleted Taxi ID: %s", taxiID)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Taxi location deleted.")
}