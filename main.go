package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
	"strings"
)

type Country struct {
	Code string
	Name string
}

// Connection to the MySQL database
var db *sql.DB

// Show a list of countries
func countryList(w http.ResponseWriter, r *http.Request) {
	// Query the database
	rows, err := db.Query("SELECT code, name FROM countries ORDER BY name")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying database: %v", err), 500)
		return
	}

	// Fetch the rows
	var countries []Country
	for rows.Next() {
		var item Country
		err = rows.Scan(&item.Code, &item.Name)
		countries = append(countries, item)
	}

	// Output the result
	t := template.Must(template.ParseFiles("tmpl/country-list.html"))
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

	rows, err := db.Query("SELECT code, name FROM countries WHERE code=?", code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying database: %v", err), 500)
		return
	}

	// Fetch the rows
	if !rows.Next() {
		http.Error(w, fmt.Sprintf("Country (code: %s) not found", code), 404)
		return
	}

	var item Country
	err = rows.Scan(&item.Code, &item.Name)

	// Output the result
	t := template.Must(template.ParseFiles("tmpl/country-view.html"))
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

		stmt, err := db.Prepare("INSERT INTO countries (code, name) VALUES(?, ?)")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error inserting country: %v", err), 500)
			return
		}
		stmt.Exec(item.Code, item.Name)
		http.Redirect(w, r, "/country/list", 302)
		return
	}

	t := template.Must(template.ParseFiles("tmpl/country-add.html"))
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

  rows, err := db.Query("SELECT code, name FROM countries WHERE code=?", original_code)
  if err != nil {
    http.Error(w, fmt.Sprintf("Error reading country: %v", err), 500)
    return
  }
  if !rows.Next() {
    http.Error(w, fmt.Sprintf("Error, country (code: %s) not found", original_code), 404)
    return
  }

  var item Country
  rows.Scan(&item.Code, &item.Name)

	t := template.Must(template.ParseFiles("tmpl/country-edit.html"))
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
	var code = parts[3]

	// Delete the row
	stmt, err := db.Prepare("DELETE FROM countries WHERE code=?")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting country: %v", err), 500)
		return
	}

	res, err := stmt.Exec(code)

	numAffected, err := res.RowsAffected()
	if numAffected == 0 {
		http.Error(w, fmt.Sprintf("Error deleting country, country (code: %s) not found", code), 404)
		return
	}

	fmt.Fprintf(w, "Delete successful")
}

func main() {
  // Connect to the MySQL database
  var err error
	db, err = sql.Open("mysql", "ian:FI0wxB@tcp(192.168.1.82)/family")
	if err != nil {
		panic(err)
	}
  defer db.Close()

  // Setup routes
	http.HandleFunc("/country/add", countryAdd)
  http.HandleFunc("/country/edit/", countryEdit)
	http.HandleFunc("/country/delete/", countryDelete)
	http.HandleFunc("/country/list", countryList)
	http.HandleFunc("/country/view/", countryView)
	http.ListenAndServe(":8090", nil)
}
