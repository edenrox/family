package main

import (
	"database/sql"
	"fmt"
	"log"
)

type Country struct {
	Code string
	Name string
}

func LoadCountryByCode(db *sql.DB, code string) (*Country, error) {
	log.Printf("Load Country code: %s", code)
	rows, err := db.Query("SELECT code, name FROM countries WHERE code=?", code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Fetch the rows
	if !rows.Next() {
		return nil, fmt.Errorf("No rows found matching code: %s", code)
	}

	var item Country
	err = rows.Scan(&item.Code, &item.Name)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func LoadCountryList(db *sql.DB) ([]Country, error) {
	log.Printf("Load Country list")
	rows, err := db.Query("SELECT code, name FROM countries ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Fetch the rows
	var countries []Country
	for rows.Next() {
		var item Country
		err = rows.Scan(&item.Code, &item.Name)
		countries = append(countries, item)
	}
	return countries, nil
}

func DeleteCountryByCode(db *sql.DB, code string) error {
	log.Printf("Delete Country code: %s", code)
	// Delete the row
	res, err := db.Exec("DELETE FROM countries WHERE code=?", code)
	if err != nil {
		return err
	}

	numAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if numAffected == 0 {
		return fmt.Errorf("Country not found, code: %s", code)
	}
	return nil
}
