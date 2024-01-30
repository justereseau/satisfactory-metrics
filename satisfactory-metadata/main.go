package main

import (
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/lib/pq"
)

// Define parameters
var (
	logLevel      = flag.String("log.level", "info", "Only log messages with the given severity or above. One of: [debug, info, warn, error, none]")
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

	// Print connection string for debugging purposes
	fmt.Println(psqlconn)

	// Connect to database
	db, err := sql.Open("postgres", psqlconn)

	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	defer db.Close()

	err = db.Ping()

	if err != nil {
		panic("Failed to ping database: " + err.Error())
	}

	// 	cacheWorker := NewCacheWorker("http://"+frmHostname+":"+strconv.Itoa(frmPort), db)
	// 	go cacheWorker.Start()

	// 	fmt.Printf(`
	// FRM Cache started
	// Press ctrl-c to exit`)

	// 	// Wait for an interrupt signal
	// 	sigChan := make(chan os.Signal, 1)
	// 	if runtime.GOOS == "windows" {
	// 		signal.Notify(sigChan, os.Interrupt)
	// 	} else {
	// 		signal.Notify(sigChan, syscall.SIGTERM)
	// 		signal.Notify(sigChan, syscall.SIGINT)
	// 	}
	// 	<-sigChan

	// 	cacheWorker.Stop()
}
