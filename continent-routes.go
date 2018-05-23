package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func continentList(w http.ResponseWriter, r *http.Request) {
	continents, err := LoadContinentList(db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading continents: %v", err), 500)
		return
	}
	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/continent/list.html")).Execute(w, continents)
	if err != nil {
		panic(err)
	}
}

func continentView(w http.ResponseWriter, r *http.Request) {
	continentCode, err := getPathParam(r, "continentCode", 3)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading continentCode: %v", err), 500)
		return
	}

	continent, err := LoadContinentByCode(db, continentCode)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading continent: %v", err), 500)
		return
	}

	countries, err := LoadCountriesByContinentCode(db, continentCode)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading countries: %v", err), 500)
		return
	}

	data := struct {
		Continent *ContinentWithMap
		Countries []Country
	}{
		continent,
		countries,
	}

	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/continent/view.html")).Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func addContinentRoutes() {
	http.HandleFunc("/continent/list", continentList)
	http.HandleFunc("/continent/view/", continentView)
}
