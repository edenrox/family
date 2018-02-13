package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func personView(w http.ResponseWriter, r *http.Request) {
	personId, err := getIntPathParam(r, "personId", 3)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not parse person id: %v", err), 400)
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
	err := template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/person/list.html")).Execute(w, nil)
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

func personJsonSearch(w http.ResponseWriter, r *http.Request) {
	prefix := r.FormValue("prefix")
	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		offset = 0
	}
	people, err := LoadPersonLiteListByNamePrefix(db, prefix, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading people: %v", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

func addPersonRoutes() {
	http.HandleFunc("/person/json/search", personJsonSearch)
	http.HandleFunc("/person/list", personList)
	http.HandleFunc("/person/calendar", personCalendar)
	http.HandleFunc("/person/view/", personView)
}
