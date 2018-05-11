package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func cityList(w http.ResponseWriter, r *http.Request) {
	err := template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/city/list.html")).Execute(w, "data")
	if err != nil {
		panic(err)
	}
}

func cityView(w http.ResponseWriter, r *http.Request) {
	cityId, err := getIntPathParam(r, "cityId", 3 /* index */)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing cityId: %v", err), 400)
		return
	}

	city, err := LoadCityById(db, cityId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading city: %v", err), 500)
		return
	}
	region, err := LoadRegionById(db, city.RegionId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading region: %v", err), 500)
		return
	}

	personList, err := LoadPersonLiteListByHomeCityId(db, cityId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading people: %v", err), 500)
		return
	}

	data := struct {
		City   *CityLite
		Region *RegionLite
		People []PersonLite
	}{
		city,
		region,
		personList,
	}

	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/city/view.html")).Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func cityAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		cityName := r.FormValue("name")
		regionId, err := strconv.Atoi(r.FormValue("region_id"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid region_id: %s", r.FormValue("region_id")), 400)
			return
		}
		lat, err := strconv.ParseFloat(r.FormValue("lat"), 32)
		lng, err := strconv.ParseFloat(r.FormValue("lng"), 32)
		_, err = InsertCity(db, cityName, regionId, float32(lat), float32(lng))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating city: %v", err), 500)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/region/view/%d", regionId), 302)
	}

	err := template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/city/add.html")).Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func cityJsonSearch(w http.ResponseWriter, r *http.Request) {
	prefix := strings.Trim(r.FormValue("prefix"), " \t")
	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		offset = 0
	}
	cities, err := LoadCitiesByPrefix(db, prefix, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading cities: %v", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cities)
}

func cityDelete(w http.ResponseWriter, r *http.Request) {
	cityId, err := getIntPathParam(r, "cityId" /* index= */, 3)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing cityId: %v", err), 400)
		return
	}

	err = DeleteCity(db, cityId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error delting city (id: %d): %v", cityId, err), 500)
		return
	}

	http.Redirect(w, r, "/country/list", 302)
}

func addCityRoutes() {
	http.HandleFunc("/city/list/", cityList)
	http.HandleFunc("/city/view/", cityView)
	http.HandleFunc("/city/json/search", cityJsonSearch)
	http.HandleFunc("/city/add", cityAdd)
	http.HandleFunc("/city/delete/", cityDelete)
}
