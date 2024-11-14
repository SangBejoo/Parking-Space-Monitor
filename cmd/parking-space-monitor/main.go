// cmd/your-app/main.go
package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/SangBejoo/parking-space-monitor/internal/handlers"
    "github.com/SangBejoo/parking-space-monitor/internal/repository"
    "github.com/SangBejoo/parking-space-monitor/internal/scheduler"
    "github.com/SangBejoo/parking-space-monitor/pkg/utils"
)

func main() {
    // Initialize DB
    connStr := "user=root dbname=subagiya1 password=secret host=localhost port=5431 sslmode=disable"
    db := utils.InitDB(connStr)
    defer db.Close()

    // Initialize repositories
    taxiRepo := &repository.TaxiRepository{DB: db}
    placeRepo := &repository.PlaceRepository{DB: db}
    mappingRepo := &repository.MappingRepository{DB: db}

    // Initialize scheduler repository
    repo := &scheduler.Repository{
        DB:          db,
        TaxiRepo:    taxiRepo,
        PlaceRepo:   placeRepo,
        MappingRepo: mappingRepo,
    }

    // Initialize scheduler
    sched := scheduler.NewScheduler(repo)

    // Initialize handlers
    taxiHandler := &handlers.TaxiHandler{Repo: taxiRepo}
    placeHandler := &handlers.PlaceHandler{Repo: placeRepo}
    mappingHandler := &handlers.MappingHandler{Scheduler: sched}

    // Initialize router
    router := mux.NewRouter()

    // Register CRUD routes for Taxi Locations
    router.HandleFunc("/taxi", taxiHandler.CreateTaxiLocation).Methods("POST")
    router.HandleFunc("/taxi", taxiHandler.GetAllTaxiLocations).Methods("GET")
    router.HandleFunc("/taxi/{id}", taxiHandler.GetTaxiLocation).Methods("GET")
    router.HandleFunc("/taxi/{id}", taxiHandler.UpdateTaxiLocation).Methods("PUT")
    router.HandleFunc("/taxi/{id}", taxiHandler.DeleteTaxiLocation).Methods("DELETE")

    // Register CRUD routes for Places
    router.HandleFunc("/place", placeHandler.CreatePlace).Methods("POST")
    router.HandleFunc("/place", placeHandler.GetAllPlaces).Methods("GET")
    router.HandleFunc("/place/{id}", placeHandler.GetPlace).Methods("GET")
    router.HandleFunc("/place/{id}", placeHandler.UpdatePlace).Methods("PUT")
    router.HandleFunc("/place/{id}", placeHandler.DeletePlace).Methods("DELETE")

    // Register routes for Mapping
	router.HandleFunc("/mapping/trigger", mappingHandler.TriggerMapping).Methods("POST")
	router.HandleFunc("/mapping", mappingHandler.GetMapping).Methods("GET")
    router.HandleFunc("/mapping/{id}", mappingHandler.GetMapping).Methods("GET")


    // Start the server
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatalf("Could not start server: %s\n", err.Error())
    }
}