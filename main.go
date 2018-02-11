package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Connection to the MySQL database
var db *sql.DB

func personView(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Error, no country specified", 400)
		return
	}
	personId, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Could not parse person id: "+parts[3], 400)
		return
	}

	person, err := LoadPersonById(db, personId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading person: %v", err), 500)
		return
	}

	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/person/view.html")).Execute(w, person)
	if err != nil {
		panic(err)
	}
}

func personList(w http.ResponseWriter, r *http.Request) {
	personList, err := LoadPersonLiteList(db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading people list: %v", err), 500)
		return
	}

	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/person/list.html")).Execute(w, personList)
	if err != nil {
		panic(err)
	}
}

func personCalendar(w http.ResponseWriter, r *http.Request) {
	calendar, err := LoadPeopleCalendar(db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading person calendar: %v", err), 500)
		return
	}

	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/person/calendar.html")).Execute(w, calendar)
	if err != nil {
		panic(err)
	}
}

func spouseDelete(w http.ResponseWriter, r *http.Request) {
	person1Id, err := strconv.Atoi(r.FormValue("person1_id"))
	person2Id, err := strconv.Atoi(r.FormValue("person2_id"))

	err = DeleteSpouse(db, person1Id, person2Id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting spouse: %v", err), 500)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/person/view/%d", person1Id), 302)
}

func spouseAdd(w http.ResponseWriter, r *http.Request) {
	person1Id, err := strconv.Atoi(r.FormValue("person1_id"))

	if r.Method == "POST" {
		person2Id, err := strconv.Atoi(r.FormValue("person2_id"))
		status, err := strconv.Atoi(r.FormValue("status"))

		err = InsertSpouse(db, person1Id, person2Id, status)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error adding spouse: %v", err), 500)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/person/view/%d", person1Id), 302)
		return
	}

	person1, err := LoadPersonLiteById(db, person1Id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading person: %v", err), 500)
		return
	}

	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/spouse/add.html")).Execute(w, person1)
	if err != nil {
		panic(err)
	}
}

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

	http.HandleFunc("/spouse/add", spouseAdd)
	http.HandleFunc("/spouse/delete", spouseDelete)
	http.HandleFunc("/person/list", personList)
	http.HandleFunc("/person/calendar", personCalendar)
	http.HandleFunc("/person/view/", personView)
	http.ListenAndServe(":8090", nil)
}
