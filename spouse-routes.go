package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

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

	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/spouse/add.html")).Execute(w, person1)
	if err != nil {
		panic(err)
	}
}

func addSpouseRoutes() {
	http.HandleFunc("/spouse/add", spouseAdd)
	http.HandleFunc("/spouse/delete", spouseDelete)
}
