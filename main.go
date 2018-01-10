package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

// Connection to the MySQL database
var db *sql.DB

// Show a list of countries
func countryList(w http.ResponseWriter, r *http.Request) {
	countries, err := LoadCountryList(db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading country list: %v", err), 500)
		return
	}

	// Output the result
	t := template.Must(template.ParseFiles("tmpl/country/list.html"))
	err = t.Execute(w, countries)
	if err != nil {
		panic(err)
	}
}

// Show a single country
func countryView(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Error, no country specified", 400)
		return
	}
	var code = parts[3]

	item, err := LoadCountryByCode(db, code)

	// Output the result
	t := template.Must(template.ParseFiles("tmpl/country/view.html"))
	err = t.Execute(w, item)
	if err != nil {
		panic(err)
	}
}

// Add a new country
func countryAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		item := Country{
			Code: r.FormValue("code"),
			Name: r.FormValue("name"),
		}
		if item.Code == "" || item.Name == "" {
			http.Error(w, "Bad request, empty code or name", 400)
			return
		}

		_, err := db.Exec("INSERT INTO countries (code, name) VALUES(?, ?)", item.Code, item.Name)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error inserting country: %v", err), 500)
			return
		}
		http.Redirect(w, r, "/country/list", 302)
		return
	}

	t := template.Must(template.ParseFiles("tmpl/country/add.html"))
	err := t.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

// Add a new country
func countryEdit(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Error, no country specified", 400)
		return
	}
	var original_code = parts[3]

	if r.Method == "POST" {
		item := Country{
			Code: r.FormValue("code"),
			Name: r.FormValue("name"),
		}
		if item.Code == "" || item.Name == "" {
			http.Error(w, "Bad request, empty code or name", 400)
			return
		}

		_, err := db.Exec("UPDATE countries SET code=?, name=? WHERE code=?", item.Code, item.Name, original_code)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating country: %v", err), 500)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/country/view/%s", item.Code), 302)
		return
	}

	item, err := LoadCountryByCode(db, original_code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading country: %v", err), 500)
		return
	}

	t := template.Must(template.ParseFiles("tmpl/country/edit.html"))
	err = t.Execute(w, item)
	if err != nil {
		panic(err)
	}
}

// Delete a country
func countryDelete(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Error, no country specified", 400)
		return
	}
	code := parts[3]

	err := DeleteCountryByCode(db, code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting country: %v", err), 500)
		return
	}
	http.Redirect(w, r, "/country/list", 302)
}

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

	t := template.Must(template.ParseFiles("tmpl/person/view.html"))
	err = t.Execute(w, person)
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

	t := template.Must(template.ParseFiles("tmpl/person/list.html"))
	err = t.Execute(w, personList)
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

	t := template.Must(template.ParseFiles("tmpl/person/calendar.html"))
	err = t.Execute(w, calendar)
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

	t := template.Must(template.ParseFiles("tmpl/spouse/add.html"))
	err = t.Execute(w, person1)
	if err != nil {
		panic(err)
	}
}

func regionList(w http.ResponseWriter, r *http.Request) {
	regions, err := LoadRegionList(db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading regions: %v", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(regions)
}

func cityAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		cityName := r.FormValue("name")
		regionId, err := strconv.Atoi(r.FormValue("region_id"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid region_id: %s", r.FormValue("region_id")), 400)
			return
		}
		_, err = InsertCity(db, cityName, regionId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating city: %v", err), 500)
			return
		}
		http.Redirect(w, r, "/city/list", 302)
	}

	t := template.Must(template.ParseFiles("tmpl/city/add.html"))
	err := t.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func citySearch(w http.ResponseWriter, r *http.Request) {
	prefix := strings.Trim(r.FormValue("prefix"), " \t")
	if len(prefix) == 0 {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "[]")
		return
	}
	cities, err := LoadCitiesByPrefix(db, prefix)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading cities: %v", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cities)
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
	//http.HandleFunc("/city/view/", cityView)
	http.HandleFunc("/city/search", citySearch)
	http.HandleFunc("/city/add", cityAdd)
	//http.HandleFunc("/city/delete/", cityDelete)
	http.HandleFunc("/region/list", regionList)
	http.HandleFunc("/spouse/add", spouseAdd)
	http.HandleFunc("/spouse/delete", spouseDelete)
	http.HandleFunc("/person/list", personList)
	http.HandleFunc("/person/calendar", personCalendar)
	http.HandleFunc("/person/view/", personView)
	http.HandleFunc("/country/add", countryAdd)
	http.HandleFunc("/country/edit/", countryEdit)
	http.HandleFunc("/country/delete/", countryDelete)
	http.HandleFunc("/country/list", countryList)
	http.HandleFunc("/country/view/", countryView)
	http.ListenAndServe(":8090", nil)
}
