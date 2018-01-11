package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type PersonLite struct {
	Id     int
	Name   string
	Gender string
}

type Person struct {
	Id         int
	FirstName  string
	MiddleName string
	LastName   string
	NickName   string
	fullName   string
	Gender     string
	IsAlive    bool
	BirthDate  time.Time
	BirthCity  *CityLite
	HomeCity   *CityLite
	Mother     *PersonLite
	Father     *PersonLite
	Children   []PersonLite
	Spouses    []SpouseLite
	Siblings   []PersonLite
}

func BuildFullName(firstName string, middleName string, lastName string, nickName string) string {
	fullName := firstName
	if nickName != "" {
		fullName += " \"" + nickName + "\""
	}
	if middleName != "" {
		fullName += " " + middleName
	}
	if lastName != "" {
		fullName += " " + lastName
	}
	return fullName
}

func (p *Person) FullName() string {
	if p.fullName == "" {
		p.fullName = BuildFullName(p.FirstName, p.MiddleName, p.LastName, p.NickName)
	}
	return p.fullName
}

func (p *Person) BirthDateFormatted() string {
	return p.BirthDate.Format("Mon, Jan 2, 2006")
}

func (p *Person) HasBirthDate() bool {
	return p.BirthDate != time.Time{}
}

func (p *Person) Age() int {
	now := time.Now()
	years := now.Year() - p.BirthDate.Year()

	if now.Month() < p.BirthDate.Month() {
		years -= 1
	} else if now.Month() == p.BirthDate.Month() && now.Day() < p.BirthDate.Day() {
		years -= 1
	}

	return years
}

func LoadPersonLiteById(db *sql.DB, id int) (*PersonLite, error) {
	log.Printf("Load PersonLite id: %d", id)
	rows, err := db.Query("SELECT id, first_name, middle_name, last_name, nick_name, gender"+
		" FROM people"+
		" WHERE id=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("Person not found with id: %d", id)
	}
	item := readPersonLiteFromRows(rows)
	return &item, nil
}

func LoadPersonById(db *sql.DB, id int) (*Person, error) {
	log.Printf("Load Person id: %d", id)
	rows, err := db.Query("SELECT id, first_name, middle_name, last_name, nick_name, mother_id, father_id, birth_date, is_alive, home_city_id, birth_city_id, gender FROM people WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		rows.Close()
		return nil, fmt.Errorf("Person not found with id: %d", id)
	}
	var item Person
	var motherId, fatherId, homeCityId, birthCityId sql.NullInt64
	var birthDateString sql.NullString
	var gender string
	err = rows.Scan(&item.Id, &item.FirstName, &item.MiddleName, &item.LastName, &item.NickName, &motherId, &fatherId, &birthDateString, &item.IsAlive, &homeCityId, &birthCityId, &gender)
	rows.Close()
	if err != nil {
		return nil, err
	}
	if motherId.Valid {
		item.Mother, err = LoadPersonLiteById(db, int(motherId.Int64))
	}
	if fatherId.Valid {
		item.Father, err = LoadPersonLiteById(db, int(fatherId.Int64))
	}
	log.Printf("birthDate: %v", birthDateString)
	if birthDateString.Valid {
		item.BirthDate, err = time.Parse("2006-01-02", birthDateString.String)
	}
	if gender == "M" {
		item.Gender = "Male"
	} else {
		item.Gender = "Female"
	}
	if birthCityId.Valid {
		item.BirthCity, err = LoadCityById(db, int(birthCityId.Int64))
		if err != nil {
			return nil, err
		}
	}
	if homeCityId.Valid {
		item.HomeCity, _ = LoadCityById(db, int(homeCityId.Int64))
		if err != nil {
			return nil, err
		}
	}
	item.Children, err = LoadChildrenPersonLite(db, id)
	if err != nil {
		return nil, err
	}
	item.Spouses, err = LoadSpousesByPersonId(db, id)
	if err != nil {
		return nil, err
	}
	item.Siblings, err = LoadSiblingsPersonLite(db, id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func LoadPersonLiteList(db *sql.DB) ([]PersonLite, error) {
	log.Printf("Load PersonLite list")
	rows, err := db.Query("SELECT id, first_name, middle_name, last_name, nick_name, gender" +
		" FROM people" +
		" ORDER BY last_name, first_name, middle_name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []PersonLite
	for rows.Next() {
		item := readPersonLiteFromRows(rows)
		list = append(list, item)
	}
	return list, nil
}

func LoadChildrenPersonLite(db *sql.DB, personId int) ([]PersonLite, error) {
	log.Printf("Load children PersonLite list, id: %d", personId)
	rows, err := db.Query("SELECT id, first_name, middle_name, last_name, nick_name, gender"+
		" FROM people "+
		" WHERE father_id = ? OR mother_id = ?", personId, personId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []PersonLite
	for rows.Next() {
		item := readPersonLiteFromRows(rows)
		list = append(list, item)
	}
	return list, nil
}

func LoadSiblingsPersonLite(db *sql.DB, personId int) ([]PersonLite, error) {
	log.Printf("Load sibling PersonLite list, id: %d", personId)
	row := db.QueryRow("SELECT father_id, mother_id FROM people WHERE id=?", personId)
	var fatherId, motherId int
	row.Scan(&fatherId, &motherId)

	rows, err := db.Query("SELECT id, first_name, middle_name, last_name, nick_name, gender"+
		" FROM people"+
		" WHERE father_id=? AND mother_id=? AND id!=?", fatherId, motherId, personId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []PersonLite
	for rows.Next() {
		item := readPersonLiteFromRows(rows)
		list = append(list, item)
	}
	return list, nil
}

func readPersonLiteFromRows(rows *sql.Rows) PersonLite {
	var id int
	var firstName, middleName, lastName, nickName, gender string
	rows.Scan(&id, &firstName, &middleName, &lastName, &nickName, &gender)
	if gender == "M" {
		gender = "Male"
	} else {
		gender = "Female"
	}
	return PersonLite{Id: id, Name: BuildFullName(firstName, middleName, lastName, nickName), Gender: gender}
}
