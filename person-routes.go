package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
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

func personGraph(w http.ResponseWriter, r *http.Request) {
	personList, err := LoadPersonLiteList(db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading people: %v", err), 500)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "digraph PersonGraph {")

	parentMap := make(map[string]bool)
	for _, person := range personList {
		// Add the person to the graph
		nameNoQuotes := strings.Replace(person.Name, "\"", "", -1)
		fmt.Fprintf(w, "    P_%d [label=\"%s\"]\n", person.Id, nameNoQuotes)

		// Add a relationship for the parents if one exists
		if (person.MotherId > 0) || (person.FatherId > 0) {
			parentKey := fmt.Sprintf("R_%d_%d", person.MotherId, person.FatherId)
			if parentMap[parentKey] == false {
				parentMap[parentKey] = true
				fmt.Fprintf(w, "    %s [label=\"\" shape=point]\n", parentKey)
				if person.MotherId > 0 {
					fmt.Fprintf(w, "    P_%d -> %s\n", person.MotherId, parentKey)
				}
				if person.FatherId > 0 {
					fmt.Fprintf(w, "    P_%d -> %s\n", person.FatherId, parentKey)
				}
			}
			fmt.Fprintf(w, "    %s -> P_%d\n", parentKey, person.Id)
		}
	}
	fmt.Fprintln(w, "}")
}

func addPersonRoutes() {
	http.HandleFunc("/person/json/search", personJsonSearch)
	http.HandleFunc("/person/list", personList)
	http.HandleFunc("/person/calendar", personCalendar)
	http.HandleFunc("/person/view/", personView)
	http.HandleFunc("/person/graph/", personGraph)
}
