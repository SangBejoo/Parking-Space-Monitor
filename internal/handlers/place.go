// internal/handlers/place.go
package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
	"fmt"

    "database/sql"
    "github.com/gorilla/mux"
    "github.com/SangBejoo/parking-space-monitor/internal/models"
    "github.com/SangBejoo/parking-space-monitor/internal/repository"
)

// PlaceHandler handles HTTP requests for Place operations.
type PlaceHandler struct {
    Repo *repository.PlaceRepository
}

// CreatePlace handles the creation of a new place.
func (ph *PlaceHandler) CreatePlace(w http.ResponseWriter, r *http.Request) {
    var place models.Place
    if err := json.NewDecoder(r.Body).Decode(&place); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    placeID, err := ph.Repo.CreatePlace(place)
    if err != nil {
        http.Error(w, "Failed to create place", http.StatusInternalServerError)
        return
    }

    place.PlaceID = placeID
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(place)
}

// GetAllPlaces retrieves all places.
func (ph *PlaceHandler) GetAllPlaces(w http.ResponseWriter, r *http.Request) {
    places, err := ph.Repo.GetAllPlaces()
    if err != nil {
        http.Error(w, "Failed to query places", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(places)
}

// GetPlace retrieves a single place by ID.
func (ph *PlaceHandler) GetPlace(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    placeID, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid place ID", http.StatusBadRequest)
        return
    }

    place, err := ph.Repo.GetPlaceByID(placeID)
    if err == sql.ErrNoRows {
        http.Error(w, "Place not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, "Failed to query place", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(place)
}

// UpdatePlace handles updating an existing place.
func (ph *PlaceHandler) UpdatePlace(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    placeID, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid place ID", http.StatusBadRequest)
        return
    }

    var place models.Place
    if err := json.NewDecoder(r.Body).Decode(&place); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if err := ph.Repo.UpdatePlace(placeID, place); err != nil {
        if err.Error() == "place not found" {
            http.Error(w, "Place not found", http.StatusNotFound)
            return
        }
        http.Error(w, "Failed to update place", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Place updated.")
}

// DeletePlace handles deleting a place by ID.
func (ph *PlaceHandler) DeletePlace(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    placeID, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid place ID", http.StatusBadRequest)
        return
    }

    if err := ph.Repo.DeletePlace(placeID); err != nil {
        if err.Error() == "place not found" {
            http.Error(w, "Place not found", http.StatusNotFound)
            return
        }
        http.Error(w, "Failed to delete place", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Place deleted.")
}