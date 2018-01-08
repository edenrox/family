package main

import (
	"database/sql"
	"log"
	"time"
)

type PeopleCalendar struct {
	Months []CalendarMonth
}

type CalendarMonth struct {
	Name   string
	People []CalendarPerson
}

type CalendarPerson struct {
	Day       int
	Person    PersonLite
	BirthDate time.Time
}

func (p *CalendarPerson) BirthDateFormatted() string {
	return p.BirthDate.Format("Mon, Jan 2, 2006")
}

func LoadPeopleCalendar(db *sql.DB) (*PeopleCalendar, error) {
	log.Printf("Loading people calendar")
	rows, err := db.Query(
		"SELECT id, first_name, middle_name, last_name, nick_name, birth_date" +
			" FROM people" +
			" WHERE is_alive = 1 AND birth_date IS NOT NULL" +
			" ORDER BY MONTH(birth_date), DAY(birth_date)")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	calendar := PeopleCalendar{}
	var lastMonth *CalendarMonth
	var item *CalendarPerson
	for rows.Next() {
		item = new(CalendarPerson)
		var id int
		var firstName, middleName, lastName, nickName string
		var birthDateString sql.NullString
		rows.Scan(&id, &firstName, &middleName, &lastName, &nickName, &birthDateString)
		birthDate, err := time.Parse("2006-01-02", birthDateString.String)
		if err != nil {
			return nil, err
		}
		birthMonth := birthDate.Format("January")
		fullName := BuildFullName(firstName, middleName, lastName, nickName)
		item.Person = PersonLite{Id: id, Name: fullName}
		item.Day = birthDate.Day()
		item.BirthDate = birthDate

		if lastMonth == nil || lastMonth.Name != birthMonth {
			if lastMonth != nil {
				calendar.Months = append(calendar.Months, *lastMonth)
			}
			lastMonth = new(CalendarMonth)
			lastMonth.Name = birthMonth
			lastMonth.People = make([]CalendarPerson, 0, 10)
		}
		lastMonth.People = append(lastMonth.People, *item)
	}
	if lastMonth != nil {
		calendar.Months = append(calendar.Months, *lastMonth)
	}
	return &calendar, nil
}
