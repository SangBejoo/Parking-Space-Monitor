// cmd/parking-space-monitor/main.go
package main

import (
    "database/sql"
    "log"
    "net/http"

    "github.com/SangBejoo/parking-space-monitor/internal/handlers"
    "github.com/SangBejoo/parking-space-monitor/internal/repository"
    "github.com/gorilla/mux"
    "github.com/SangBejoo/parking-space-monitor/internal/scheduler"
    _ "github.com/lib/pq"
)

func main() {
    router := mux.NewRouter()

    // Initialize database connection
    db, err := sql.Open("postgres", "user=root dbname=subagiya1 password=secret host=localhost port=5431 sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create tables if they do not exist
    err = CreateTablesIfNotExist(db)
    if err != nil {
        log.Fatal("Failed to create tables: ", err)
    }

    // Initialize repositories
    taxiRepo := &repository.TaxiRepository{DB: db}
    placeRepo := &repository.PlaceRepository{DB: db}
    mappingRepo := &repository.MappingRepository{DB: db}

    // Initialize handlers with repositories
    taxiHandler := &handlers.TaxiHandler{Repo: taxiRepo}
    placeHandler := &handlers.PlaceHandler{Repo: placeRepo}

    // Initialize scheduler and mapping handler
    schedulerRepo := &scheduler.Repository{
        DB:          db,
        TaxiRepo:    taxiRepo,
        PlaceRepo:   placeRepo,
        MappingRepo: mappingRepo,
    }
    sched := scheduler.NewScheduler(schedulerRepo)
    mappingHandler := &handlers.MappingHandler{Scheduler: sched}

    // Define the bulk insert route before the dynamic {id} routes
    router.HandleFunc("/taxis/bulk", taxiHandler.CreateMultipleTaxiLocations).Methods("POST")

    router.HandleFunc("/taxis", taxiHandler.CreateMultipleTaxiLocations).Methods("POST")
    router.HandleFunc("/taxis", taxiHandler.GetAllTaxiLocations).Methods("GET")
    router.HandleFunc("/taxis/{id}", taxiHandler.GetTaxiLocation).Methods("GET")
    router.HandleFunc("/taxis/{id}", taxiHandler.UpdateTaxiLocation).Methods("PUT")
    router.HandleFunc("/taxis/{id}", taxiHandler.DeleteTaxiLocation).Methods("DELETE")

    router.HandleFunc("/places/bulk", placeHandler.CreateMultiplePlaces).Methods("POST")
    router.HandleFunc("/places", placeHandler.CreatePlace).Methods("POST")
    router.HandleFunc("/places", placeHandler.GetAllPlaces).Methods("GET")
    router.HandleFunc("/places/{id}", placeHandler.GetPlace).Methods("GET")
    router.HandleFunc("/places/{id}", placeHandler.UpdatePlace).Methods("PUT")
    router.HandleFunc("/places/{id}", placeHandler.DeletePlace).Methods("DELETE")

    // Define mapping routes
    router.HandleFunc("/mapping/trigger", mappingHandler.TriggerMapping).Methods("POST")
    router.HandleFunc("/mapping", mappingHandler.GetMapping).Methods("GET")

    http.Handle("/", router)
    log.Println("Server is running on port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}

// CreateTablesIfNotExist creates necessary tables if they do not already exist.
func CreateTablesIfNotExist(db *sql.DB) error {
    tableCreationQueries := []string{
        `CREATE TABLE IF NOT EXISTS taxi_location (
            taxi_id VARCHAR PRIMARY KEY,
            longitude FLOAT NOT NULL,
            latitude FLOAT NOT NULL,
            updated_at TIMESTAMP NOT NULL
        );`,
        `CREATE TABLE IF NOT EXISTS places (
            place_id SERIAL PRIMARY KEY,
            place_name VARCHAR NOT NULL,
            polygon JSONB NOT NULL
        );`,
        `CREATE TABLE IF NOT EXISTS mapping (
            taxi_id VARCHAR,
            place_id INT,
            PRIMARY KEY (taxi_id, place_id),
            FOREIGN KEY (taxi_id) REFERENCES taxi_location(taxi_id),
            FOREIGN KEY (place_id) REFERENCES places(place_id)
        );`,
        `CREATE TABLE IF NOT EXISTS counters (
            taxi_id VARCHAR,
            place_id INT,
            counter INT DEFAULT 0,
            last_counted TIMESTAMP,
            PRIMARY KEY (taxi_id, place_id),
            FOREIGN KEY (taxi_id) REFERENCES taxi_location(taxi_id),
            FOREIGN KEY (place_id) REFERENCES places(place_id)
        );`,
    }

    for _, query := range tableCreationQueries {
        _, err := db.Exec(query)
        if err != nil {
            return err
        }
    }

    log.Println("All necessary tables are ensured to exist.")
    return nil
}