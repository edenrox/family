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

type CountryFull struct {
	Code        string
	Name        string
	CapitalCity *CityLite
	Gdp         int
	Population  int
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

func LoadFullCountryList(db *sql.DB) ([]CountryFull, error) {
	defer trace(traceName("LoadFullCountryList"))
	rows, err := db.Query(
		"SELECT co.code, co.name, co.gdp, co.population, ci.city_id, ci.city_name, ci.region_id, ci.region_code" +
			" FROM countries co" +
			"   LEFT JOIN city_view ci ON ci.city_id = co.capital_city_id" +
			" ORDER BY co.name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Fetch the rows
	var countries []CountryFull
	for rows.Next() {
		var item CountryFull
		var cityId, cityRegionId sql.NullInt64
		var cityName, cityRegionCode sql.NullString
		err = rows.Scan(&item.Code, &item.Name, &item.Gdp, &item.Population, &cityId, &cityName, &cityRegionId, &cityRegionCode)
		if err != nil {
			return nil, err
		}
		if cityId.Valid {
			item.CapitalCity = &CityLite{
				Id:         int(cityId.Int64),
				Name:       cityName.String,
				RegionId:   int(cityRegionId.Int64),
				RegionAbbr: cityRegionCode.String,
			}
		}

		countries = append(countries, item)
	}
	return countries, nil
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
