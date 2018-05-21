package main

import (
	"database/sql"
	"fmt"
	"time"
)

type Holiday struct {
	Id   int
	Date time.Time
	Name string
}

type YearHolidays struct {
	Year     int
	Holidays []Holiday
}

func LoadHolidaysByYear(db *sql.DB, startYear int) ([]YearHolidays, error) {
	defer trace(traceName(fmt.Sprintf("LoadHolidaysByYear(%d)", startYear)))

	holidays, err := LoadHolidays(db, startYear)
	if err != nil {
		return nil, err
	}
	var result []YearHolidays
	var year *YearHolidays
	for _, value := range holidays {
		if year == nil || year.Year != value.Date.Year() {
			result = append(result, YearHolidays{})
			year = &result[len(result)-1]
			year.Year = value.Date.Year()
		}
		year.Holidays = append(year.Holidays, value)
	}
	return result, nil
}

func LoadHolidays(db *sql.DB, startYear int) ([]Holiday, error) {
	defer trace(traceName(fmt.Sprintf("LoadHolidays(%d)", startYear)))
	rows, err := db.Query(
		"SELECT id, date, name"+
			" FROM holidays "+
			" WHERE date >=?"+
			" ORDER BY date ASC", fmt.Sprintf("%04d-01-01", startYear))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var holidays []Holiday
	for rows.Next() {
		var holiday Holiday
		var holidayDateString string
		err := rows.Scan(&holiday.Id, &holidayDateString, &holiday.Name)
		if err != nil {
			return nil, err
		}
		holiday.Date, _ = time.Parse("2006-01-02", holidayDateString)
		holidays = append(holidays, holiday)
	}
	return holidays, nil
}
