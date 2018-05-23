package main

import (
	"database/sql"
	"fmt"
	"math"
)

type CityLite struct {
	Id          int
	Name        string
	RegionId    int
	RegionAbbr  string
	CountryAbbr string
	Latitude    float64
	Longitude   float64
}

func (c *CityLite) Format() string {
	return fmt.Sprintf("%s, %s, %s", c.Name, c.RegionAbbr, c.CountryAbbr)
}

func (c *CityLite) HasLocation() bool {
	return c.Latitude != 0
}

func (c *CityLite) FormatLocation() string {
	latDir := "N"
	if c.Latitude < 0 {
		latDir = "S"
	}
	longDir := "E"
	if c.Longitude < 0 {
		longDir = "W"
	}
	return fmt.Sprintf("%.3f %s, %.3f %s", math.Abs(c.Latitude), latDir, math.Abs(c.Longitude), longDir)
}

func (c *CityLite) MapUrl() string {
	return fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/staticmap"+
			"?center=%.4f,%.4f"+
			"&zoom=5"+
			"&size=500x300"+
			"&maptype=roadmap"+
			"&markers=color:blue|%.4f,%.4f"+
			"&key=%s",
		c.Latitude, c.Longitude,
		c.Latitude, c.Longitude,
		config.mapsApiKey)
}

func readCityListFromRows(rows *sql.Rows) ([]CityLite, error) {
	var cities []CityLite
	for rows.Next() {
		city, err := readCityFromRows(rows)
		if err != nil {
			return nil, err
		}
		cities = append(cities, *city)
	}
	return cities, nil
}

func readCityFromRows(rows *sql.Rows) (*CityLite, error) {
	city := CityLite{}
	err := rows.Scan(&city.Id, &city.Name, &city.RegionId, &city.RegionAbbr, &city.CountryAbbr, &city.Latitude, &city.Longitude)
	if err != nil {
		return nil, err
	}
	return &city, nil
}

func LoadCityById(db *sql.DB, id int) (*CityLite, error) {
	defer trace(traceName(fmt.Sprintf("LoadCityById(%d)", id)))
	rows, err := db.Query(
		"SELECT city_id, city_name, region_id, region_code, country_code, lat, lng"+
			" FROM city_view "+
			" WHERE city_id=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, fmt.Errorf("Error, city not found (id: %d)", id)
	}

	city, err := readCityFromRows(rows)
	if err != nil {
		return nil, err
	}
	return city, nil
}

func LoadCitiesByRegionId(db *sql.DB, regionId int) ([]CityLite, error) {
	defer trace(traceName(fmt.Sprintf("LoadCitiesByRegionId(%d)", regionId)))
	rows, err := db.Query(
		"SELECT city_id, city_name, region_id, region_code, country_code, lat, lng"+
			" FROM city_view"+
			" WHERE region_id=?"+
			" ORDER BY city_name", regionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return readCityListFromRows(rows)
}

func LoadCitiesByCountryCode(db *sql.DB, countryCode string) ([]CityLite, error) {
	defer trace(traceName(fmt.Sprintf("LoadCitiesByCountryCode(%s)", countryCode)))
	rows, err := db.Query(
		"SELECT city_id, city_name, region_id, region_code, country_code, lat, lng"+
			" FROM city_view"+
			" WHERE country_code=?"+
			" ORDER BY region_code, city_name", countryCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return readCityListFromRows(rows)
}

func LoadCitiesByPrefix(db *sql.DB, prefix string, offset int) ([]CityLite, error) {
	defer trace(traceName(fmt.Sprintf("LoadCitiesByPrefix(%s)", prefix)))
	prefix = prefix + "%"
	rows, err := db.Query(
		"SELECT city_id, city_name, region_id, region_code, country_code, lat, lng"+
			" FROM city_view"+
			" WHERE city_name LIKE ?"+
			" ORDER BY city_name, region_code, country_code"+
			" LIMIT ?, 10", prefix, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return readCityListFromRows(rows)
}

func DeleteCity(db *sql.DB, id int) error {
	trace(traceName(fmt.Sprintf("DeleteCity(%d)", id)))
	res, err := db.Exec("DELETE FROM cities WHERE id=?", id)
	if err != nil {
		return err
	}
	numAffected, _ := res.RowsAffected()
	if numAffected < 1 {
		return fmt.Errorf("City not found (id: %d)", id)
	}
	return nil
}

func InsertCity(db *sql.DB, name string, regionId int, lat float32, lng float32) (*CityLite, error) {
	trace(traceName(fmt.Sprintf("InsertCity(name:%s, regionId: %d)", name, regionId)))
	res, err := db.Exec("INSERT INTO cities (name, region_id, lat, lng) VALUES(?, ?, ?, ?)", name, regionId, lat, lng)
	if err != nil {
		return nil, err
	}
	cityId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return LoadCityById(db, int(cityId))
}

func UpdateCity(db *sql.DB, cityId int, name string, regionId int, lat float32, lng float32) error {
	trace(traceName(fmt.Sprintf("UpdateCity(%d)", cityId)))
	_, err := db.Exec(
		"UPDATE cities"+
			" SET name=?, region_id=?, lat=?, lng=?"+
			" WHERE id=?",
		name, regionId, lat, lng, cityId)
	return err
}
