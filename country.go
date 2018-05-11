package main

import (
	"database/sql"
	"fmt"
	"log"
)

type Country struct {
	Code          string
	Name          string
	CapitalCityId int
}

func LoadCountryByCode(db *sql.DB, code string) (*Country, error) {
	defer trace(traceName(fmt.Sprintf("LoadCountryBycode(%s)", code)))
	rows, err := db.Query(
		"SELECT code, name, capital_city_id"+
			" FROM countries"+
			" WHERE code=?", code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Fetch the rows
	if !rows.Next() {
		return nil, fmt.Errorf("No rows found matching code: %s", code)
	}

	item, err := readCountryFromRows(rows)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func LoadCountryList(db *sql.DB) ([]Country, error) {
	defer trace(traceName("LoadCountryList"))
	rows, err := db.Query(
		"SELECT code, name, capital_city_id" +
			" FROM countries" +
			" ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Fetch the rows
	var countries []Country
	for rows.Next() {
		item, err := readCountryFromRows(rows)
		if err != nil {
			return nil, err
		}
		countries = append(countries, *item)
	}
	return countries, nil
}

func readCountryFromRows(rows *sql.Rows) (*Country, error) {
	country := Country{}
	var capitalCityId sql.NullInt64
	err := rows.Scan(&country.Code, &country.Name, &capitalCityId)
	if capitalCityId.Valid {
		country.CapitalCityId = int(capitalCityId.Int64)
	}
	if err != nil {
		return nil, err
	}
	return &country, nil
}

func DeleteCountryByCode(db *sql.DB, code string) error {
	defer trace(traceName(fmt.Sprintf("DeleteCountryByCode(%s)", code)))
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

func InsertCountry(db *sql.DB, name string, code string, capitalCityId int) error {
	log.Printf("Insert country (name: %s, code: %s, capitalCityId: %d)", name, code, capitalCityId)
	nullableCapitalCityId := sql.NullInt64{
		Int64: int64(capitalCityId),
		Valid: capitalCityId > 0,
	}
	res, err := db.Exec("INSERT INTO countries (name, code, capital_city_id) VALUES(?, ?, ?)", name, code, nullableCapitalCityId)
	if err != nil {
		return err
	}
	_, err = res.LastInsertId()
	return err
}

func UpdateCountry(db *sql.DB, originalCode string, name string, code string, capitalCityId int) error {
	log.Printf("Update country (originalCode: %s, name: %s, code: %s, capitalCityId: %d)", originalCode, name, code, capitalCityId)
	nullableCapitalCityId := sql.NullInt64{
		Int64: int64(capitalCityId),
		Valid: capitalCityId > 0,
	}
	_, err := db.Exec("UPDATE countries SET name=?, code=?, capital_city_id=? WHERE code=?", name, code, nullableCapitalCityId, originalCode)
	return err
}
