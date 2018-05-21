package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"sort"
	"time"
)

type PeopleCalendar struct {
	Months []CalendarMonth
}

type CalendarMonth struct {
	Name   string
	Events []CalendarEvent
}

type CalendarEvent struct {
	Date    time.Time
	Type    string
	Caption template.HTML
}

type CalendarPerson struct {
	BirthDate time.Time
	Person    PersonLite
}

type CalendarAnniversary struct {
	MarriedDate time.Time
	Person1     PersonLite
	Person2     PersonLite
}

func LoadPeopleCalendar(db *sql.DB) (*PeopleCalendar, error) {
	defer trace(traceName("LoadPeopleCalendar"))

	personLookup, err := loadPeopleByBirthMonth(db)
	if err != nil {
		return nil, err
	}
	anniversaryLookup, err := loadAnniversariesByMonth(db)
	if err != nil {
		return nil, err
	}

	calendar := PeopleCalendar{
		Months: make([]CalendarMonth, 12),
	}

	personTemplate := template.Must(template.New("person").Parse("<a href=\"/person/view/{{.Person.Id}}\">{{.Person.Name}}</a>"))
	anniversaryTemplate := template.Must(template.New("anniversary").Parse(
		"<a href=\"/person/view/{{.Person1.Id}}\">{{.Person1.Name}}</a> &amp;" +
			" <a href=\"/person/view/{{.Person2.Id}}\">{{.Person2.Name}}</a>"))

	for i := 0; i < 12; i++ {
		calendar.Months[i] = CalendarMonth{
			Name: time.Month(i + 1).String(),
		}

		events := make([]CalendarEvent, 0)

		for _, value := range (*personLookup)[i] {
			var buf bytes.Buffer
			personTemplate.Execute(&buf, value)
			event := CalendarEvent{
				Date:    value.BirthDate,
				Type:    "Birthday",
				Caption: template.HTML(buf.String()),
			}
			events = append(events, event)
		}

		for _, value := range (*anniversaryLookup)[i] {
			var buf bytes.Buffer
			anniversaryTemplate.Execute(&buf, value)
			event := CalendarEvent{
				Date:    value.MarriedDate,
				Type:    "Anniversary",
				Caption: template.HTML(buf.String()),
			}
			events = append(events, event)
		}
		sort.Slice(events, func(i, j int) bool { return events[i].Date.Day() < events[j].Date.Day() })
		calendar.Months[i].Events = events
	}

	return &calendar, nil
}

func loadPeopleByBirthMonth(db *sql.DB) (*[][]CalendarPerson, error) {
	defer trace(traceName("LoadPeopleByBirthMonth"))
	monthLookup := make([][]CalendarPerson, 12)
	for i := 0; i < 12; i++ {
		monthLookup[i] = make([]CalendarPerson, 0, 0)
	}

	rows, err := db.Query(
		"SELECT id, first_name, middle_name, last_name, nick_name, gender, birth_date" +
			" FROM people" +
			" WHERE is_alive = 1 AND birth_date IS NOT NULL" +
			" ORDER BY MONTH(birth_date), DAY(birth_date)")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var birthDateString, firstName, middleName, lastName, nickName, gender string
		rows.Scan(&id, &firstName, &middleName, &lastName, &nickName, &gender, &birthDateString)
		birthDate, err := time.Parse("2006-01-02", birthDateString)
		if err != nil {
			fmt.Printf("Error parsing birthdate; person id: %d, birthdate: %s\n", id, birthDateString)
			return nil, err
		}

		item := CalendarPerson{
			BirthDate: birthDate,
			Person: PersonLite{
				Id:     id,
				Name:   BuildFullName(firstName, middleName, lastName, nickName),
				Gender: GetGenderName(gender),
			},
		}

		monthIndex := birthDate.Month() - 1
		monthLookup[monthIndex] = append(monthLookup[monthIndex], item)
		count++
	}

	log.Printf("Birthdays found: %d", count)
	return &monthLookup, nil
}

func loadAnniversariesByMonth(db *sql.DB) (*[][]CalendarAnniversary, error) {
	defer trace(traceName("LoadAnniversariesByMonth"))
	monthLookup := make([][]CalendarAnniversary, 12)
	for i := 0; i < 12; i++ {
		monthLookup[i] = make([]CalendarAnniversary, 0, 0)
	}

	rows, err := db.Query(
		"SELECT s.married_date, " +
			"  p1.id, p1.first_name, p1.middle_name, p1.last_name, p1.nick_name, p1.gender, " +
			"  p2.id, p2.first_name, p2.middle_name, p2.last_name, p2.nick_name, p2.gender " +
			" FROM spouses s" +
			"  INNER JOIN people p1 ON p1.id = s.person1_id" +
			"  INNER JOIN people p2 ON p2.id = s.person2_id" +
			" WHERE s.status = 1 AND s.married_date IS NOT NULL" +
			" ORDER BY MONTH(s.married_date), DAY(s.married_date)")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var marriedDateString string
		var p1Id, p2Id int
		var p1FirstName, p1MiddleName, p1LastName, p1NickName, p1Gender string
		var p2FirstName, p2MiddleName, p2LastName, p2NickName, p2Gender string
		rows.Scan(&marriedDateString,
			&p1Id, &p1FirstName, &p1MiddleName, &p1LastName, &p1NickName, &p1Gender,
			&p2Id, &p2FirstName, &p2MiddleName, &p2LastName, &p2NickName, &p2Gender)

		marriedDate, err := time.Parse("2006-01-02", marriedDateString)
		if err != nil {
			return nil, err
		}

		item := CalendarAnniversary{
			MarriedDate: marriedDate,
			Person1: PersonLite{
				Id:     p1Id,
				Name:   BuildFullName(p1FirstName, p1MiddleName, p1LastName, p1NickName),
				Gender: GetGenderName(p1Gender),
			},
			Person2: PersonLite{
				Id:     p2Id,
				Name:   BuildFullName(p2FirstName, p2MiddleName, p2LastName, p2NickName),
				Gender: GetGenderName(p2Gender),
			},
		}

		monthIndex := marriedDate.Month() - 1
		monthLookup[monthIndex] = append(monthLookup[monthIndex], item)
		count++
	}

	log.Printf("Anniversaries found: %d", count)
	return &monthLookup, nil
}
