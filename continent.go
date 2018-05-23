package main

import (
	"database/sql"
	"fmt"
)

type Continent struct {
	Code  string
	Name  string
	Color string
}

type ContinentWithMap struct {
	Code         string
	Name         string
	MapLatitude  float32
	MapLongitude float32
	MapZoom      int
	Color        string
}

func (c *ContinentWithMap) MapUrl() string {
	return fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/staticmap"+
			"?center=%.4f,%.4f"+
			"&zoom=%d"+
			"&size=500x400"+
			"&maptype=roadmap"+
			"&key=%s",
		c.MapLatitude, c.MapLongitude,
		c.MapZoom,
		config.mapsApiKey)
}

func LoadContinentList(db *sql.DB) ([]Continent, error) {
	defer trace(traceName("LoadContinentList"))
	rows, err := db.Query(
		"SELECT code, name, color" +
			" FROM continents" +
			" ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return readContinentListFromRows(rows)
}

func LoadContinentByCode(db *sql.DB, code string) (*ContinentWithMap, error) {
	defer trace(traceName(fmt.Sprintf("LoadContinentByCode(%s)", code)))
	rows, err := db.Query(
		"SELECT code, name, map_lat, map_lng, map_zoom, color"+
			" FROM continents"+
			" WHERE code=?",
		code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("Error continent not found. Code: %s", code)
	}
	var item ContinentWithMap
	err = rows.Scan(&item.Code, &item.Name, &item.MapLatitude, &item.MapLongitude, &item.MapZoom, &item.Color)
	if err != nil {
		return nil, fmt.Errorf("Error scanning continent: %v", err)
	}
	return &item, nil
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
	err := rows.Scan(&continent.Code, &continent.Name, &continent.Color)
	if err != nil {
		return nil, fmt.Errorf("Error scanning continent: %v", err)
	}
	return &continent, nil
}
