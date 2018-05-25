package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type GrandParents struct {
	GrandFather *PersonLite
	GrandMother *PersonLite
}

func personAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var data PersonData
		data.FirstName = strings.TrimSpace(r.FormValue("first_name"))
		data.MiddleName = strings.TrimSpace(r.FormValue("middle_name"))
		data.LastName = strings.TrimSpace(r.FormValue("last_name"))
		data.NickName = strings.TrimSpace(r.FormValue("nick_name"))

		if data.FirstName == "" {
			http.Error(w, "Error, first name can not be empty.", 400)
			return
		}

		if r.FormValue("gender") == "male" {
			data.Gender = "M"
		} else {
			data.Gender = "F"
		}
		data.BirthDate = r.FormValue("birth_date")
		if r.FormValue("is_birth_year_guess") == "1" {
			data.IsBirthYearGuess = true
		}
		if r.FormValue("is_alive") == "1" {
			data.IsAlive = true
		}
		data.BirthCityId, _ = strconv.Atoi(r.FormValue("birth_city_id"))
		data.HomeCityId, _ = strconv.Atoi(r.FormValue("home_city_id"))
		data.MotherId, _ = strconv.Atoi(r.FormValue("mother_id"))
		data.FatherId, _ = strconv.Atoi(r.FormValue("father_id"))

		personId, err := InsertPerson(db, data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error inserting person: %v", err), 500)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/person/view/%d", personId), 302)
		return
	}

	err := template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/person/add.html")).Execute(w, "empty data")
	if err != nil {
		panic(err)
	}
}

func personEdit(w http.ResponseWriter, r *http.Request) {
	personId, err := getIntPathParam(r, "personId", 3)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not parse person id: %v", err), 400)
		return
	}

	if r.Method == "POST" {
		var data PersonData
		data.FirstName = strings.TrimSpace(r.FormValue("first_name"))
		data.MiddleName = strings.TrimSpace(r.FormValue("middle_name"))
		data.LastName = strings.TrimSpace(r.FormValue("last_name"))
		data.NickName = strings.TrimSpace(r.FormValue("nick_name"))

		if data.FirstName == "" {
			http.Error(w, "Error, first name can not be empty.", 400)
			return
		}

		if r.FormValue("gender") == "male" {
			data.Gender = "M"
		} else {
			data.Gender = "F"
		}
		data.BirthDate = r.FormValue("birth_date")
		if r.FormValue("is_birth_year_guess") == "1" {
			data.IsBirthYearGuess = true
		}
		if r.FormValue("is_alive") == "1" {
			data.IsAlive = true
		}
		data.BirthCityId, _ = strconv.Atoi(r.FormValue("birth_city_id"))
		data.HomeCityId, _ = strconv.Atoi(r.FormValue("home_city_id"))
		data.MotherId, _ = strconv.Atoi(r.FormValue("mother_id"))
		data.FatherId, _ = strconv.Atoi(r.FormValue("father_id"))

		err = UpdatePerson(db, personId, data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating person: %v", err), 500)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/person/view/%d", personId), 302)
		return
	}

	person, err := LoadPersonById(db, personId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading person: %v", err), 500)
		return
	}

	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/person/edit.html")).Execute(w, person)
	if err != nil {
		panic(err)
	}
}

func LoadPersonLiteOrNull(db *sql.DB, personId int) (*PersonLite, error) {
	if personId == 0 {
		return nil, nil
	} else {
		return LoadPersonLiteById(db, personId)
	}
}

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

	var paternalGrandParents, maternalGrandParents GrandParents
	if person.Father != nil {
		paternalGrandParents.GrandMother, _ = LoadPersonLiteOrNull(db, person.Father.MotherId)
		paternalGrandParents.GrandFather, _ = LoadPersonLiteOrNull(db, person.Father.FatherId)
	}
	if person.Mother != nil {
		maternalGrandParents.GrandMother, _ = LoadPersonLiteOrNull(db, person.Mother.MotherId)
		maternalGrandParents.GrandFather, _ = LoadPersonLiteOrNull(db, person.Mother.FatherId)
	}

	data := struct {
		Person               *Person
		PaternalGrandParents GrandParents
		MaternalGrandParents GrandParents
	}{
		person,
		paternalGrandParents,
		maternalGrandParents,
	}

	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/person/view.html", "tmpl/person/tree.html")).Execute(w, data)
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

func personDelete(w http.ResponseWriter, r *http.Request) {
	personId, err := getIntPathParam(r, "personId", 3)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not parse person id: %v", err), 400)
		return
	}
	err = DeletePerson(db, personId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete person: %v", err), 500)
		return
	}
	http.Redirect(w, r, "/person/list", 302)
}

func addPersonRoutes() {
	http.HandleFunc("/person/json/search", personJsonSearch)
	http.HandleFunc("/person/list", personList)
	http.HandleFunc("/person/calendar", personCalendar)
	http.HandleFunc("/person/view/", personView)
	http.HandleFunc("/person/add", personAdd)
	http.HandleFunc("/person/delete/", personDelete)
	http.HandleFunc("/person/edit/", personEdit)
	http.HandleFunc("/person/graph/", personGraph)
}
