package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

// Show a list of countries
func countryList(w http.ResponseWriter, r *http.Request) {
	countries, err := LoadCountryList(db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading country list: %v", err), 500)
		return
	}

	// Output the result
	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/country/list.html")).Execute(w, countries)
	if err != nil {
		panic(err)
	}
}

// List all countries as JSON
func countryJsonList(w http.ResponseWriter, r *http.Request) {
	countries, err := LoadCountryList(db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading country list: %v", err), 500)
		return
	}

	// Output the result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(countries)
}

// Show a single country
func countryView(w http.ResponseWriter, r *http.Request) {
	code, err := getPathParam(r, "countryCode", 3)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing countryCode: %v", err), 400)
	}

	item, err := LoadCountryByCode(db, code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading country: %v", err), 500)
		return
	}
	regions, err := LoadRegionsByCountryCode(db, code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading regions: %v", err), 500)
		return
	}
	data := struct {
		Country *Country
		Regions []RegionLite
	}{
		item,
		regions,
	}

	// Output the result
	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/country/view.html")).Execute(w, data)
	if err != nil {
		panic(err)
	}
}

// Add a new country
func countryAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		capitalCityId, _ := strconv.Atoi(r.FormValue("capital_city_id"))
		gdp, _ := strconv.Atoi(r.FormValue("gdp"))
		population, _ := strconv.Atoi(r.FormValue("population"))
		item := CountryData{
			Code:           strings.TrimSpace(r.FormValue("code")),
			Name:           strings.TrimSpace(r.FormValue("name")),
			CapitalCityId:  capitalCityId,
			Gdp:            gdp,
			Population:     population,
			HasRegionIcons: r.FormValue("has_region_icons") == "1",
		}
		if item.Code == "" || item.Name == "" {
			http.Error(w, "Bad request, empty code or name", 400)
			return
		}

		err := InsertCountry(db, item)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error inserting country: %v", err), 500)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/country/view/%s", item.Code), 302)
		return
	}

	err := template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/country/add.html")).Execute(w, nil)
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
	var originalCode = parts[3]

	if r.Method == "POST" {
		capitalCityId, _ := strconv.Atoi(r.FormValue("capital_city_id"))
		gdp, _ := strconv.Atoi(r.FormValue("gdp"))
		population, _ := strconv.Atoi(r.FormValue("population"))
		item := CountryData{
			Code:           strings.TrimSpace(r.FormValue("code")),
			Name:           strings.TrimSpace(r.FormValue("name")),
			CapitalCityId:  capitalCityId,
			Gdp:            gdp,
			Population:     population,
			HasRegionIcons: r.FormValue("has_region_icons") == "1",
		}
		if item.Code == "" || item.Name == "" {
			http.Error(w, "Bad request, empty code or name", 400)
			return
		}

		err := UpdateCountry(db, originalCode, item)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating country: %v", err), 500)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/country/view/%s", item.Code), 302)
		return
	}

	item, err := LoadCountryByCode(db, originalCode)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading country: %v", err), 500)
		return
	}

	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/country/edit.html")).Execute(w, item)
	if err != nil {
		panic(err)
	}
}

// Delete a country
func countryDelete(w http.ResponseWriter, r *http.Request) {
	code, err := getPathParam(r, "countryCode", 3)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing countryCode: %v", err), 400)
	}

	err = DeleteCountryByCode(db, code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting country: %v", err), 500)
		return
	}
	http.Redirect(w, r, "/country/list", 302)
}

func addCountryRoutes() {
	http.HandleFunc("/country/add", countryAdd)
	http.HandleFunc("/country/edit/", countryEdit)
	http.HandleFunc("/country/delete/", countryDelete)
	http.HandleFunc("/country/list", countryList)
	http.HandleFunc("/country/json/list", countryJsonList)
	http.HandleFunc("/country/view/", countryView)
}
