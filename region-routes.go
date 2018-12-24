package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type RegionGroup struct {
	CountryName string
	CountryCode string
	Regions     []RegionLite
}

func regionList(w http.ResponseWriter, r *http.Request) {
	regions, err := LoadRegionList(db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading regions: %v", err), 500)
		return
	}

	var regionGroups []RegionGroup
	var group *RegionGroup
	for _, region := range regions {
		if group == nil || group.CountryCode != region.CountryCode {
			item := RegionGroup{
				CountryName: region.CountryName,
				CountryCode: region.CountryCode,
			}
			regionGroups = append(regionGroups, item)
			group = &regionGroups[len(regionGroups)-1]
		}
		group.Regions = append(group.Regions, region)
	}

	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/region/list.html")).Execute(w, regionGroups)
	if err != nil {
		panic(err)
	}
}

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
		data := RegionData{
			Name:        strings.TrimSpace(r.FormValue("name")),
			Code:        strings.TrimSpace(r.FormValue("code")),
			CountryCode: strings.TrimSpace(r.FormValue("country_code")),
		}
		if data.Name == "" || data.CountryCode == "" {
			http.Error(w, fmt.Sprintf("Error, empty name or country code: %v", data), 400)
			return
		}
		_, err := InsertRegion(db, data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating region: %v", err), 500)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/country/view/%s", data.CountryCode), 302)
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

func regionEdit(w http.ResponseWriter, r *http.Request) {
	regionId, err := getIntPathParam(r, "regionId", 3)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing regionId: %v", err), 400)
		return
	}

	if r.Method == "POST" {
		data := RegionData{
			Name:        strings.TrimSpace(r.FormValue("name")),
			Code:        strings.TrimSpace(r.FormValue("code")),
			CountryCode: strings.TrimSpace(r.FormValue("country_code")),
		}
		if data.Name == "" || data.CountryCode == "" {
			http.Error(w, fmt.Sprintf("Error, empty name or country code: %v", data), 400)
			return
		}
		err = UpdateRegion(db, regionId, data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating region: %v", err), 500)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/country/view/%s", data.CountryCode), 302)
		return
	}

	region, err := LoadRegionById(db, regionId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading region: %v", err), 400)
		return
	}

	countries, err := LoadCountryList(db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading countries: %v", err), 500)
		return
	}

	data := struct {
		Countries []Country
		Region    *RegionLite
	}{
		countries,
		region,
	}

	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/region/edit.html")).Execute(w, data)
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

	http.Redirect(w, r, fmt.Sprintf("/country/view/%s", region.CountryCode), 302)
}

func addRegionRoutes() {
	http.HandleFunc("/region/add", regionAdd)
	http.HandleFunc("/region/edit/", regionEdit)
	http.HandleFunc("/region/list", regionList)
	http.HandleFunc("/region/json/list", regionJsonList)
	http.HandleFunc("/region/view/", regionView)
	http.HandleFunc("/region/delete/", regionDelete)
}
