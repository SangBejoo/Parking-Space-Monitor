// internal/handlers/taxi.go
package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
    "database/sql"
    "github.com/SangBejoo/parking-space-monitor/internal/models"
    "github.com/SangBejoo/parking-space-monitor/internal/repository"
)

// TaxiHandler handles HTTP requests for Taxi operations.
type TaxiHandler struct {
    Repo *repository.TaxiRepository
}

// CreateTaxiLocation handles the creation of a new taxi location.
func (th *TaxiHandler) CreateTaxiLocation(w http.ResponseWriter, r *http.Request) {
    var location models.TaxiLocation
    if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if err := th.Repo.CreateTaxi(location); err != nil {
        http.Error(w, "Failed to create taxi location", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
    fmt.Fprintf(w, "Taxi location created.")
}

// GetAllTaxiLocations retrieves all taxi locations.
func (th *TaxiHandler) GetAllTaxiLocations(w http.ResponseWriter, r *http.Request) {
    taxis, err := th.Repo.GetAllTaxis()
    if err != nil {
        http.Error(w, "Failed to query taxi locations", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(taxis)
}

// GetTaxiLocation retrieves a single taxi location by ID.
func (th *TaxiHandler) GetTaxiLocation(w http.ResponseWriter, r *http.Request) {
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

// UpdateTaxiLocation handles updating an existing taxi location.
func (th *TaxiHandler) UpdateTaxiLocation(w http.ResponseWriter, r *http.Request) {
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

// DeleteTaxiLocation handles deleting a taxi location by ID.
func (th *TaxiHandler) DeleteTaxiLocation(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    taxiID := vars["id"]

    if err := th.Repo.DeleteTaxi(taxiID); err != nil {
        if err.Error() == "taxi not found" {
            http.Error(w, "Taxi not found", http.StatusNotFound)
            return
        }
        http.Error(w, "Failed to delete taxi location", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Taxi location deleted.")
}