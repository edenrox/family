package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	_ "github.com/go-sql-driver/mysql"
)

var config struct {
	databaseConnectionString string
	mapsApiKey               string
}

// Connection to the MySQL database
var db *sql.DB

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func loadFlags() {
	flag.StringVar(&config.databaseConnectionString, "database", "", "dsn for connecting to a mysql database")
	flag.StringVar(&config.mapsApiKey, "mapsApiKey", "", "API Key for connecting to Google Static Maps API")
	flag.Parse()
}

func main() {
	loadFlags()

	// Connect to the MySQL database
	fmt.Printf("Connecting to database: %s\n", config.databaseConnectionString)
	var err error
	db, err = sql.Open("mysql", config.databaseConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Setup static finals
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

	// Setup routes
	http.HandleFunc("/health", health)
	addCountryRoutes()
	addRegionRoutes()
	addCityRoutes()
	addPersonRoutes()
	addSpouseRoutes()
	addCronRoutes()
	addHolidayRoutes()
	addContinentRoutes()

	http.ListenAndServe(":8090", nil)
}
