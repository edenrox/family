package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

func holidayList(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	holidays, err := LoadHolidaysByYear(db, now.Year())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading holidays: %v", err), 500)
		return
	}

	// Output the result
	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/holiday/list.html")).Execute(w, holidays)
	if err != nil {
		panic(err)
	}
}

func addHolidayRoutes() {
	http.HandleFunc("/holiday/list", holidayList)
}
