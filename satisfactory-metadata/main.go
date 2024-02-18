package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

// Define parameters
var (
	frmApiAddress = flag.String("frm.listen-address", "http://localhost:8080", "Address of Ficsit Remote Monitoring webserver")

	pgHost     = flag.String("db.pghost", "postgres", "postgres hostname")
	pgPort     = flag.Int("db.pgport", 5432, "postgres port")
	pgPassword = flag.String("db.pgpassword", "secretpassword", "postgres password")
	pgUser     = flag.String("db.pguser", "postgres", "postgres username")
	pgDb       = flag.String("db.pgdb", "postgres", "postgres db")
)

func main() {
	// Get parameters
	flag.Parse()

	// Generate connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", *pgHost, *pgPort, *pgUser, *pgPassword, *pgDb)

	// Connect to database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	defer db.Close()

	// Ping database to ensure connection is still alive
	err = db.Ping()
	if err != nil {
		panic("Failed to ping database: " + err.Error())
	}

	// Initialize database
	err = initDB(db)
	if err != nil {
		panic("Failed to initialize database: " + err.Error())
	}

	pullMetrics(db, "factory", "/getFactory", false)
	pullMetrics(db, "extractor", "/getExtractor", false)
	pullMetrics(db, "dropPod", "/getDropPod", false)
	pullMetrics(db, "storageInv", "/getStorageInv", false)
	pullMetrics(db, "worldInv", "/getWorldInv", false)
	pullMetrics(db, "droneStation", "/getDroneStation", false)
	pullMetrics(db, "generators", "/getCoalGenerator", false)
	pullMetrics(db, "generators", "/getBiomassGenerator", false)
	pullMetrics(db, "generators", "/getFuelGenerator", false)
	pullMetrics(db, "generators", "/getNuclearGenerator", false)
	pullMetrics(db, "generators", "/getGeothermalGenerator", false)
	pullMetrics(db, "drone", "/getDrone", false)
	pullMetrics(db, "train", "/getTrains", false)
	pullMetrics(db, "truck", "/getVehicles", false)
	pullMetrics(db, "trainStation", "/getTrainStation", false)
	pullMetrics(db, "truckStation", "/getTruckStation", false)
}

// initialize the database
func initDB(db *sql.DB) error {
	_, err := db.Exec(`
	  CREATE TABLE IF NOT EXISTS cache(
		id serial primary key,
		metric text NOT NULL,
		frm_data jsonb
	  );
	  CREATE INDEX IF NOT EXISTS cache_metric_idx ON cache(metric);
	  FLUSH cache;
	  `)
	if err != nil {
		fmt.Println("Error while creating DB Table : ", err)
		return err
	}
	fmt.Println("DB Tables created successfully")
	return err
}

// pull metrics from the Ficsit Remote Monitoring API
func pullMetrics(db *sql.DB, metric string, route string, keepHistory bool) {
	resp, err := http.Get(*frmApiAddress + route)
	if err != nil {
		fmt.Println("Error while querying "+route, err)
		return
	}

	var content []json.RawMessage
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&content)
	if err != nil {
		// Try to found if it is an empty object
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return
		}
		fmt.Println("Error while decoding "+route+" response: ", err)
		return
	}
	defer resp.Body.Close()

	data := []string{}
	for _, c := range content {
		data = append(data, string(c[:]))
	}

	err = cacheMetrics(db, metric, data)

	if err != nil {
		fmt.Println("error when caching metadatas %s", err)
		return
	}

	fmt.Println("Successfully cached metadatas from " + route + " for " + metric)
}

// cache metrics in the database
func cacheMetrics(db *sql.DB, metric string, data []string) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	for _, s := range data {
		insert := `insert into cache (metric, frm_data) values($1, $2)`
		_, err = tx.Exec(insert, metric, s)
		if err != nil {
			return
		}
	}
	return
}
