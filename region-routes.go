package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func regionJsonList(w http.ResponseWriter, r *http.Request) {
	regions, err := LoadRegionList(db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading regions: %v", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(regions)
}

func regionView(w http.ResponseWriter, r *http.Request) {
	regionId, err := getIntPathParam(r, "regionId", 3)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing regionId: %v", err), 400)
		return
	}

	region, err := LoadRegionById(db, regionId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading region: %v", err), 500)
		return
	}
	cities, err := LoadCitiesByRegionId(db, regionId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading cities: %v", err), 500)
		return
	}

	data := struct {
		Region *RegionLite
		Cities []CityLite
	}{
		region,
		cities,
	}

	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/region/view.html")).Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func regionAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		regionName := r.FormValue("name")
		regionCode := r.FormValue("code")
		countryCode := r.FormValue("country_code")
		_, err := InsertRegion(db, regionName, regionCode, countryCode)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating city: %v", err), 500)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/country/view/%s", countryCode), 302)
		return
	}

	countries, err := LoadCountryList(db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading countries: %v", err), 500)
		return
	}

	data := struct {
		Countries           []Country
		SelectedCountryCode string
	}{
		countries,
		r.FormValue("country_code"),
	}

	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/region/add.html")).Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func regionDelete(w http.ResponseWriter, r *http.Request) {
	regionId, err := getIntPathParam(r, "regionId", 3)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing regionId: %v", err), 400)
		return
	}

	region, err := LoadRegionById(db, regionId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading region: %v", err), 400)
		return
	}

	err = DeleteRegion(db, regionId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting region: %v", err), 500)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/country/view/%d", region.CountryCode), 302)
}

func addRegionRoutes() {
	http.HandleFunc("/region/add", regionAdd)
	http.HandleFunc("/region/json/list", regionJsonList)
	http.HandleFunc("/region/view/", regionView)
	http.HandleFunc("/region/delete/", regionDelete)
}
