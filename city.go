package main

import (
	"database/sql"
	"fmt"
	"log"
)

type CityLite struct {
	Id          int
	Name        string
	RegionAbbr  string
	CountryAbbr string
}

func (c *CityLite) Format() string {
	return fmt.Sprintf("%s, %s, %s", c.Name, c.RegionAbbr, c.CountryAbbr)
}

func LoadCityById(db *sql.DB, id int) (*CityLite, error) {
	log.Printf("Load CityLite id: %d", id)
	rows, err := db.Query("SELECT city_name, region_code, country_code FROM city_view WHERE city_id=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, fmt.Errorf("Error, city not found (id: %d)", id)
	}

	city := CityLite{Id: id}
	err = rows.Scan(&city.Name, &city.RegionAbbr, &city.CountryAbbr)
	if err != nil {
		return nil, err
	}
	return &city, nil
}

func LoadCitiesByPrefix(db *sql.DB, prefix string) ([]CityLite, error) {
	log.Printf("Load Cities by prefix:: %s", prefix)
	prefix = prefix + "%"
	rows, err := db.Query("SELECT city_id, city_name, region_code, country_code FROM city_view WHERE city_name LIKE ?", prefix)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []CityLite
	for rows.Next() {
		city := CityLite{}
		err = rows.Scan(&city.Id, &city.Name, &city.RegionAbbr, &city.CountryAbbr)
		if err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}
	return cities, nil
}

func DeleteCity(db *sql.DB, id int) error {
	log.Printf("Delete city id: %d", id)
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

func InsertCity(db *sql.DB, name string, regionId int) (*CityLite, error) {
	log.Printf("Insert city (name: %s, regionId: %d)", name, regionId)
	res, err := db.Exec("INSERT INTO cities (name, region_id) VALUES(?, ?)", name, regionId)
	if err != nil {
		return nil, err
	}
	cityId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return LoadCityById(db, int(cityId))
}
