package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

// Define parameters
var (
	frmApiAddress = flag.String("frm.listen-address", "http://localhost:8080", "Address of Ficsit Remote Monitoring webserver")

	pgHost      = flag.String("db.pghost", "postgres", "postgres hostname")
	pgPort      = flag.Int("db.pgport", 5432, "postgres port")
	pgPassword  = flag.String("db.pgpassword", "secretpassword", "postgres password")
	pgUser      = flag.String("db.pguser", "postgres", "postgres username")
	pgDb        = flag.String("db.pgdb", "postgres", "postgres db")
	metricsFile = flag.String("metrics.file", "", "Configuration file for metrics to pull from Ficsit Remote Monitoring API")
)

// New type for metrics to pull from Ficsit Remote Monitoring API
type metric struct {
	name  string
	route string
}

func main() {
	// Get parameters
	flag.Parse()

	metrics := []metric{}

	// Read metrics from file if the file is provided
	if *metricsFile != "" {
		// Open the file
		csvfile, err := os.Open(*metricsFile)
		if err != nil {
			log.Fatalln("Couldn't open the csv file", err)
		}
		r := csv.NewReader(csvfile)

		// Read the file
		records, err := r.ReadAll()
		if err != nil {
			log.Fatalln("Couldn't read the csv file", err)
		}

		// Loop through lines
		for _, record := range records {
			if len(record) != 2 {
				log.Fatalln("Invalid line in csv file", record)
			}
			metrics = append(metrics, metric{name: record[0], route: record[1]})
		}
	} else {
		// If no file is provided, use default metrics
		metrics = append(metrics, metric{name: "factory", route: "getFactory"})
		metrics = append(metrics, metric{name: "extractor", route: "getExtractor"})
		metrics = append(metrics, metric{name: "dropPod", route: "getDropPod"})
		metrics = append(metrics, metric{name: "storageInv", route: "getStorageInv"})
		metrics = append(metrics, metric{name: "worldInv", route: "getWorldInv"})
		metrics = append(metrics, metric{name: "droneStation", route: "getDroneStation"})
		metrics = append(metrics, metric{name: "generators", route: "getGenerator"})
		metrics = append(metrics, metric{name: "drone", route: "getDrone"})
		metrics = append(metrics, metric{name: "train", route: "getTrains"})
		metrics = append(metrics, metric{name: "truck", route: "getVehicles"})
		metrics = append(metrics, metric{name: "trainStation", route: "getTrainStation"})
		metrics = append(metrics, metric{name: "truckStation", route: "getTruckStation"})
	}

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

	// Pull metrics from Ficsit Remote Monitoring API
	for _, m := range metrics {
		pullMetrics(db, m.name, m.route)
	}
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
	  TRUNCATE TABLE cache;
	  `)
	if err != nil {
		fmt.Println("Error while creating DB Table : ", err)
		return err
	}
	fmt.Println("DB Tables created successfully")
	return err
}

// pull metrics from the Ficsit Remote Monitoring API
func pullMetrics(db *sql.DB, metric string, route string) {
	resp, err := http.Get(*frmApiAddress + "/" + route)
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
