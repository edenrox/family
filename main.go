package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	_ "github.com/go-sql-driver/mysql"
)

// Connection to the MySQL database
var db *sql.DB

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func main() {
	// Connect to the MySQL database
	dsn := flag.String("database", "", "dsn for connecting to a mysql database")

	flag.Parse()

	fmt.Printf("Connecting to database: %s\n", *dsn)
	var err error
	db, err = sql.Open("mysql", *dsn)
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
