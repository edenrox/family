package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

func getNextMonday() time.Time {
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if now.Weekday() == time.Sunday {
		return now.Add(time.Hour * 24 * 1)
	}
	return now.Add(time.Hour * 24 * time.Duration(8-now.Weekday()))
}

func cronReminders(w http.ResponseWriter, r *http.Request) {
	var startTime time.Time
	startTime, err := time.Parse("2006-01-02", r.FormValue("start_date"))
	if err != nil {
		startTime = getNextMonday()
	}
	endTime := startTime.Add(time.Hour * 24 * 7)

	people, err := loadPeopleWithBirthday(db, startTime, endTime)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading people: %v", err), 500)
		return
	}

	var events []CalendarEvent
	for _, value := range people {
		event := CalendarEvent{
			Day:           value.BirthDate.Day(),
			Date:          value.BirthDate,
			DateFormatted: calendarDateFormatted(value.BirthDate),
			Type:          "Birthday",
			Caption:       template.HTML(value.Person.Name),
		}
		events = append(events, event)
	}

	data := struct {
		Events             []CalendarEvent
		StartDateFormatted string
		EndDateFormatted   string
	}{
		Events:             events,
		StartDateFormatted: startTime.Format("2006-01-02"),
		EndDateFormatted:   endTime.Format("2006-01-02"),
	}

	// Output the result
	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/cron/reminders.html")).Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func loadPeopleWithBirthday(db *sql.DB, startTime time.Time, endTime time.Time) ([]CalendarPerson, error) {
	defer trace(traceName(fmt.Sprintf("loadPeopleWithBirthday(%v, %v)", startTime, endTime)))

	rows, err := db.Query(
		"SELECT id, first_name, middle_name, last_name, nick_name, gender, birth_date"+
			" FROM people"+
			" WHERE is_alive = 1 AND birth_date IS NOT NULL"+
			" AND MONTH(birth_date) >= ? AND DAY(birth_date) >= ?"+
			" AND MONTH(birth_date) <= ? AND DAY(birth_date) < ?"+
			" ORDER BY MONTH(birth_date), DAY(birth_date)",
		startTime.Month(), startTime.Day(),
		endTime.Month(), endTime.Day())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people []CalendarPerson
	for rows.Next() {
		var birthDateString, gender, firstName, middleName, lastName, nickName string
		var id int

		rows.Scan(&id, &firstName, &middleName, &lastName, &nickName, &gender, &birthDateString)

		birthDate, err := time.Parse("2006-01-02", birthDateString)
		if err != nil {
			return nil, err
		}

		item := CalendarPerson{
			BirthDate: birthDate,
			Person: PersonLite{
				Id:     id,
				Name:   BuildFullName(firstName, middleName, lastName, nickName),
				Gender: GetGenderName(gender)},
		}

		people = append(people, item)
	}

	return people, nil
}

func addCronRoutes() {
	http.HandleFunc("/cron/reminders", cronReminders)
}
