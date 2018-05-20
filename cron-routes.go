package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"sort"
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
	// End time is 2 weeks after the start time
	endTime := startTime.Add(time.Hour * 24 * 14)

	people, err := loadPeopleWithBirthday(db, startTime, endTime)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading people: %v", err), 500)
		return
	}

	anniversaries, err := loadAnniversariesInRange(db, startTime, endTime)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading anniversaries: %v", err), 500)
		return
	}

	var events []CalendarEvent

	// Add birthdays to event calendar
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

	// Add anniversaries to event calendar
	for _, value := range anniversaries {
		event := CalendarEvent{
			Day:           value.MarriedDate.Day(),
			Date:          value.MarriedDate,
			DateFormatted: calendarDateFormatted(value.MarriedDate),
			Type:          "Anniversary",
			Caption:       template.HTML(value.Person1.Name + " &amp; " + value.Person2.Name),
		}
		events = append(events, event)
	}

	// sort the event calendar
	sort.Slice(events, func(i, j int) bool { return events[i].Day < events[j].Day })

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

func loadAnniversariesInRange(db *sql.DB, startTime time.Time, endTime time.Time) ([]CalendarAnniversary, error) {
	defer trace(traceName(fmt.Sprintf("loadAnniversariesInRange(%v, %v)", startTime, endTime)))

	rows, err := db.Query(
		"SELECT s.married_date,"+
			" p1.id, p1.first_name, p1.middle_name, p1.last_name, p1.nick_name, p1.gender, "+
			" p2.id, p2.first_name, p2.middle_name, p2.last_name, p2.nick_name, p2.gender "+
			" FROM spouses s"+
			" INNER JOIN people p1 ON p1.id = s.person1_id"+
			" INNER JOIN people p2 ON p2.id = s.person2_id"+
			" WHERE s.status = 1 AND s.married_date IS NOT NULL"+
			" AND MONTH(s.married_date) >= ? AND DAY(s.married_date) >= ?"+
			" AND MONTH(s.married_date) <= ? AND DAY(s.married_date) < ?"+
			" ORDER BY MONTH(s.married_date), DAY(s.married_date)",
		startTime.Month(), startTime.Day(),
		endTime.Month(), endTime.Day())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var anniversaries []CalendarAnniversary
	for rows.Next() {

		var person1Id, person2Id int
		var marriedDateString string
		var person1FirstName, person1MiddleName, person1LastName, person1NickName, person1Gender string
		var person2FirstName, person2MiddleName, person2LastName, person2NickName, person2Gender string

		err := rows.Scan(&marriedDateString,
			&person1Id, &person1FirstName, &person1MiddleName, &person1LastName, &person1NickName, &person1Gender,
			&person2Id, &person2FirstName, &person2MiddleName, &person2LastName, &person2NickName, &person2Gender)
		if err != nil {
			return nil, err
		}

		marriedDate, _ := time.Parse("2006-01-02", marriedDateString)
		item := CalendarAnniversary{
			MarriedDate: marriedDate,
			Person1: PersonLite{
				Id:     person1Id,
				Name:   BuildFullName(person1FirstName, person1MiddleName, person1LastName, person1NickName),
				Gender: GetGenderName(person1Gender),
			},
			Person2: PersonLite{
				Id:     person2Id,
				Name:   BuildFullName(person2FirstName, person2MiddleName, person2LastName, person2NickName),
				Gender: GetGenderName(person2Gender),
			},
		}

		anniversaries = append(anniversaries, item)
	}

	return anniversaries, nil
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
