package main

import (
	"database/sql"
	"fmt"
)

type SpouseLite struct {
	Person1 PersonLite
	Person2 PersonLite
	Status  int // 1=MARRIED,2=DATING,3=EX-MARRIED
}

func (s *SpouseLite) StatusFormatted() string {
	switch s.Status {
	case 1:
		return "Married"
	case 2:
		return "Dating"
	case 3:
		return "Ex-Married"
	default:
		return ""
	}
}

func LoadSpousesByPersonId(db *sql.DB, personId int) ([]SpouseLite, error) {
	trace(traceName(fmt.Sprintf("LoadSpousesByPersonId(%d)", personId)))
	rows, err := db.Query("SELECT person1_id, person2_id, status FROM spouses WHERE person1_id=? OR person2_id=? ORDER BY status", personId, personId)
	if err != nil {
		return nil, err
	}
	person1, err := LoadPersonLiteById(db, personId)
	if err != nil {
		return nil, err
	}
	var spouseList []SpouseLite
	for rows.Next() {
		item := SpouseLite{Person1: *person1}
		var person1Id, person2Id int
		err = rows.Scan(&person1Id, &person2Id, &item.Status)
		if err != nil {
			return nil, err
		}
		if person2Id == personId {
			person2Id = person1Id
		}
		person2, err := LoadPersonLiteById(db, person2Id)
		if err != nil {
			return nil, err
		}
		item.Person2 = *person2
		spouseList = append(spouseList, item)
	}
	return spouseList, nil
}

func DeleteSpouse(db *sql.DB, person1Id int, person2Id int) error {
	if person2Id < person1Id {
		tmp := person1Id
		person1Id = person2Id
		person2Id = tmp
	}
	trace(traceName(fmt.Sprintf("DeleteSpouse(%d, %d)", person1Id, person2Id)))
	res, err := db.Exec("DELETE FROM spouses WHERE person1_id=? AND person2_id=?", person1Id, person2Id)
	if err != nil {
		return err
	}
	numAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if numAffected < 1 {
		return fmt.Errorf("Spouse not found with person1_id: %d, person2_id: %d", person1Id, person2Id)
	}
	return nil
}

func InsertSpouse(db *sql.DB, person1Id int, person2Id int, status int, marriedDate string) error {
	if person2Id < person1Id {
		tmp := person1Id
		person1Id = person2Id
		person2Id = tmp
	}
	trace(traceName(fmt.Sprintf("InsertSpouse(%d, %d, %d, %s)", person1Id, person2Id, status, marriedDate)))
	nullableMariedDate := sql.NullString{
		String: marriedDate,
		Valid:  marriedDate != "",
	}
	_, err := db.Exec(
		"INSERT INTO spouses"+
			" (person1_id, person2_id, status, marriedDate)"+
			" VALUES(?, ?, ?, ?)",
		person1Id, person2Id, status, nullableMariedDate)
	if err != nil {
		return err
	}
	return nil
}
