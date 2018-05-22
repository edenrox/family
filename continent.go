package main

import (
	"database/sql"
	"fmt"
)

type Continent struct {
	Code string
	Name string
}

func LoadContinentList(db *sql.DB) ([]Continent, error) {
	defer trace(traceName("LoadContinentList"))
	rows, err := db.Query(
		"SELECT code, name" +
			" FROM continents" +
			" ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return readContinentListFromRows(rows)
}

func LoadContinentByCode(db *sql.DB, code string) (*Continent, error) {
	defer trace(traceName(fmt.Sprintf("LoadContinentByCode(%s)", code)))
	rows, err := db.Query(
		"SELECT code, name"+
			" FROM continents"+
			" WHERE code=?",
		code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return readContinentFromRows(rows)
	} else {
		return nil, fmt.Errorf("Error continent not found. Code: %s", code)
	}
}

func readContinentListFromRows(rows *sql.Rows) ([]Continent, error) {
	var list []Continent
	for rows.Next() {
		item, err := readContinentFromRows(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *item)
	}
	return list, nil
}

func readContinentFromRows(rows *sql.Rows) (*Continent, error) {
	var continent Continent
	err := rows.Scan(&continent.Code, &continent.Name)
	if err != nil {
		return nil, fmt.Errorf("Error scanning row: %v", err)
	}
	return &continent, nil
}
