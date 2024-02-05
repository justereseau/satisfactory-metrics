package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/benbjohnson/clock"
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

	Clock = clock.New()
	now   = Clock.Now()
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

	err = flushMetricHistory(db)
	if err != nil {
		panic("Failed to flush metric history: " + err.Error())
	}

	// Low cadence metrics
	pullMetrics(db, "factory", "/getFactory", true)
	pullMetrics(db, "extractor", "/getExtractor", true)
	pullMetrics(db, "dropPod", "/getDropPod", false)
	pullMetrics(db, "storageInv", "/getStorageInv", false)
	pullMetrics(db, "worldInv", "/getWorldInv", false)
	pullMetrics(db, "droneStation", "/getDroneStation", false)
	pullMetrics(db, "generators", "/getCoalGenerator", false)
	pullMetrics(db, "generators", "/getBiomassGenerator", false)
	pullMetrics(db, "generators", "/getFuelGenerator", false)
	pullMetrics(db, "generators", "/getNuclearGenerator", false)
	pullMetrics(db, "generators", "/getGeothermalGenerator", false)

	// Realtime metrics
	pullMetrics(db, "drone", "/getDrone", true)
	pullMetrics(db, "train", "/getTrains", true)
	pullMetrics(db, "truck", "/getVehicles", true)
	pullMetrics(db, "trainStation", "/getTrainStation", true)
	pullMetrics(db, "truckStation", "/getTruckStation", true)

}

// initialize the database
func initDB(db *sql.DB) error {
	_, err := db.Exec(`
	  CREATE TABLE IF NOT EXISTS cache(
		id serial primary key,
		metric text NOT NULL,
		frm_data jsonb
	  );

	  CREATE TABLE IF NOT EXISTS cache_with_history(
		id serial primary key,
		metric text NOT NULL,
		frm_data jsonb,
		time timestamp
	  );

	  CREATE INDEX IF NOT EXISTS cache_metric_idx ON cache(metric);
	  CREATE INDEX IF NOT EXISTS cache_with_history_metric_idx ON cache_with_history(metric);
	  `)
	if err != nil {
		fmt.Println("Error while creating DB Tables : ", err)
		return err
	}
	fmt.Println("DB Tables created successfully")
	return err
}

// flush the metric history cache
func flushMetricHistory(db *sql.DB) error {
	delete := `truncate cache_with_history;`
	_, err := db.Exec(delete)
	if err != nil {
		fmt.Println("flush metrics history db error: ", err)
	}
	return err
}

// pull metrics from the Ficsit Remote Monitoring API
func pullMetrics(db *sql.DB, metric string, route string, keepHistory bool) {
	resp, err := http.Get(*frmApiAddress + route)

	if err != nil {
		fmt.Println("Error while querying Ficsit Remote Monitoring API: ", err)
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
		fmt.Println("Error while decoding Ficsit Remote Monitoring API response: ", err)
		return
	}
	defer resp.Body.Close()

	data := []string{}
	for _, c := range content {
		data = append(data, string(c[:]))
	}

	err = cacheMetrics(db, metric, data)

	if err != nil {
		fmt.Println("error when caching metrics %s", err)
		return
	}

	if keepHistory {
		err = cacheMetricsWithHistory(db, metric, data)
		if err != nil {
			fmt.Println("error when caching metrics history %s", err)
			return
		}
	}
	fmt.Println("Successfully cached metrics for " + metric)
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

	delete := `delete from cache where metric = $1;`
	_, err = tx.Exec(delete, metric)
	if err != nil {
		return
	}
	for _, s := range data {
		insert := `insert into cache (metric, frm_data) values($1, $2)`
		_, err = tx.Exec(insert, metric, s)
		if err != nil {
			return
		}
	}
	return
}

// cache metrics in the database with history
func cacheMetricsWithHistory(db *sql.DB, metric string, data []string) (err error) {
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
		insert := `insert into cache_with_history (metric, frm_data, time) values($1,$2,$3)`
		_, err = tx.Exec(insert, metric, s, now)
		if err != nil {
			return
		}
	}

	//720 = 1 hour, 5 second increments. retain that many rows for every data.
	keep := 720 * len(data)

	delete := `delete from cache_with_history where
metric = $1 and
id NOT IN (
select id from "cache_with_history" where metric = $1
order by id desc
limit $2
);`
	_, err = tx.Exec(delete, metric, keep)
	return
}
