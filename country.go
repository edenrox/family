package main

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Country struct {
	Code           string
	Name           string
	CapitalCity    *CityLite
	Gdp            int
	Population     int
	HasRegionIcons bool
	Continent      *Continent
}

type CountryData struct {
	Code           string
	Name           string
	CapitalCityId  int
	Gdp            int
	Population     int
	HasRegionIcons bool
	ContinentCode  string
}

func (c *Country) GdpFormatted() string {
	return message.NewPrinter(language.English).Sprint(c.Gdp)
}

func (c *Country) PopulationFormatted() string {
	return message.NewPrinter(language.English).Sprint(c.Population)
}

func LoadCountriesByContinentCode(db *sql.DB, continentCode string) ([]Country, error) {
	defer trace(traceName(fmt.Sprintf("LoadCountriesByContinentCode(%s)", continentCode)))
	rows, err := db.Query(
		"SELECT co.code, co.name, co.gdp, co.population, co.has_region_icons,"+
			" ci.city_id, ci.city_name, ci.region_id, ci.region_code,"+
			" ct.code, ct.name"+
			" FROM countries co"+
			"   INNER JOIN continents ct ON ct.code = co.continent_code"+
			"   LEFT JOIN city_view ci ON ci.city_id = co.capital_city_id"+
			" WHERE co.continent_code=?", continentCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return readCountryListFromRows(rows)
}

func LoadCountryByCode(db *sql.DB, code string) (*Country, error) {
	defer trace(traceName(fmt.Sprintf("LoadCountryBycode(%s)", code)))
	rows, err := db.Query(
		"SELECT co.code, co.name, co.gdp, co.population, co.has_region_icons,"+
			" ci.city_id, ci.city_name, ci.region_id, ci.region_code,"+
			" ct.code, ct.name"+
			" FROM countries co"+
			"   INNER JOIN continents ct ON ct.code = co.continent_code"+
			"   LEFT JOIN city_view ci ON ci.city_id = co.capital_city_id"+
			" WHERE co.code=?", code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return readCountryFromRows(rows)
	} else {
		return nil, fmt.Errorf("Country not found. Code: %s", code)
	}
}

func LoadCountryList(db *sql.DB) ([]Country, error) {
	defer trace(traceName("LoadCountryList"))
	rows, err := db.Query(
		"SELECT co.code, co.name, co.gdp, co.population, co.has_region_icons," +
			" ci.city_id, ci.city_name, ci.region_id, ci.region_code," +
			" ct.code, ct.name" +
			" FROM countries co" +
			"   INNER JOIN continents ct ON ct.code = co.continent_code" +
			"   LEFT JOIN city_view ci ON ci.city_id = co.capital_city_id" +
			" ORDER BY co.name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return readCountryListFromRows(rows)
}

func readCountryListFromRows(rows *sql.Rows) ([]Country, error) {
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
	country := Country{
		Continent: &Continent{},
	}
	var capitalCityId, capitalCityRegionId sql.NullInt64
	var capitalCityName, capitalCityRegionCode sql.NullString
	err := rows.Scan(&country.Code, &country.Name, &country.Gdp, &country.Population, &country.HasRegionIcons,
		&capitalCityId, &capitalCityName, &capitalCityRegionId, &capitalCityRegionCode,
		&country.Continent.Code, &country.Continent.Name)
	if capitalCityId.Valid {
		country.CapitalCity = &CityLite{
			Id:         int(capitalCityId.Int64),
			Name:       capitalCityName.String,
			RegionId:   int(capitalCityRegionId.Int64),
			RegionAbbr: capitalCityRegionCode.String,
		}
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

func InsertCountry(db *sql.DB, data CountryData) error {
	log.Printf("Insert country (data: %v)", data)
	nullableCapitalCityId := sql.NullInt64{
		Int64: int64(data.CapitalCityId),
		Valid: data.CapitalCityId > 0,
	}
	res, err := db.Exec("INSERT INTO countries"+
		" (name, code, capital_city_id, gdp, population, has_region_icons, continent_code)"+
		" VALUES(?, ?, ?, ?, ?, ?, ?)",
		data.Name, data.Code, nullableCapitalCityId, data.Gdp, data.Population, data.HasRegionIcons, data.ContinentCode)
	if err != nil {
		return err
	}
	_, err = res.LastInsertId()
	return err
}

func UpdateCountry(db *sql.DB, originalCode string, data CountryData) error {
	log.Printf("Update country (originalCode: %s, data: %v)", originalCode, data)
	nullableCapitalCityId := sql.NullInt64{
		Int64: int64(data.CapitalCityId),
		Valid: data.CapitalCityId > 0,
	}
	_, err := db.Exec("UPDATE countries "+
		"SET name=?, code=?, capital_city_id=?, gdp=?, population=?, has_region_icons=?, continent_code=? "+
		"WHERE code=?",
		data.Name, data.Code, nullableCapitalCityId, data.Gdp, data.Population, data.HasRegionIcons, data.ContinentCode,
		originalCode)
	return err
}
