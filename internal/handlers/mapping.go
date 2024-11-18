// internal/handlers/mapping.go
package handlers

import (
    "encoding/json"
    "fmt"  
    "net/http"
    "strconv"

    "database/sql" // Added import

    "github.com/SangBejoo/parking-space-monitor/internal/models"
    "github.com/SangBejoo/parking-space-monitor/internal/repository"
    "github.com/SangBejoo/parking-space-monitor/internal/scheduler"
    "github.com/gorilla/mux"
)

// MappingHandler handles HTTP requests for Mapping operations.
type MappingHandler struct {
    Repo *repository.MappingRepository
    Scheduler *scheduler.Scheduler
}

// CreateMapping handles the creation of a new mapping.
func (mh *MappingHandler) CreateMapping(w http.ResponseWriter, r *http.Request) {
    var mapping models.Mapping
    if err := json.NewDecoder(r.Body).Decode(&mapping); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if err := mh.Repo.InsertMapping(mapping); err != nil {
        http.Error(w, "Failed to create mapping", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(mapping)
}

// GetAllMappings retrieves all mappings.
func (mh *MappingHandler) GetAllMappings(w http.ResponseWriter, r *http.Request) {
    mappings, err := mh.Repo.GetAllMappings()
    if err != nil {
        http.Error(w, "Failed to retrieve mappings", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(mappings)
}

// GetMapping retrieves a single mapping by ID.
func (mh *MappingHandler) GetMapping(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    mappingID, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid mapping ID", http.StatusBadRequest)
        return
    }

    mapping, err := mh.Repo.GetMappingByID(mappingID)
    if err == sql.ErrNoRows {
        http.Error(w, "Mapping not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, "Failed to retrieve mapping", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(mapping)
}

// UpdateMapping handles updating an existing mapping.
func (mh *MappingHandler) UpdateMapping(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    mappingID, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid mapping ID", http.StatusBadRequest)
        return
    }

    var mapping models.Mapping
    if err := json.NewDecoder(r.Body).Decode(&mapping); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    mapping.ID = mappingID
    if err := mh.Repo.UpdateMapping(mapping); err != nil {
        if err.Error() == "mapping not found" {
            http.Error(w, "Mapping not found", http.StatusNotFound)
            return
        }
        http.Error(w, "Failed to update mapping", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(mapping)
}

// DeleteMapping handles deleting a mapping by ID.
func (mh *MappingHandler) DeleteMapping(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    mappingID, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid mapping ID", http.StatusBadRequest)
        return
    }

    if err := mh.Repo.DeleteMapping(mappingID); err != nil {
        if err.Error() == "mapping not found" {
            http.Error(w, "Mapping not found", http.StatusNotFound)
            return
        }
        http.Error(w, "Failed to delete mapping", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Mapping deleted successfully"))
}

// Additional Handlers for Scheduler Integration

// TriggerMapping manually triggers the mapTaxiLocations function.
func (mh *MappingHandler) TriggerMapping(w http.ResponseWriter, r *http.Request) {
    go func() {
        mh.Scheduler.MapTaxiLocations()
    }()
    fmt.Fprintf(w, "Mapping process triggered manually.")
}