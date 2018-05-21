package main

import (
	"database/sql"
	"fmt"
)

type RegionLite struct {
	Id            int
	Code          string
	Name          string
	CountryCode   string
	CountryName   string
	HasRegionIcon bool
}

type RegionData struct {
	Code        string
	Name        string
	CountryCode string
}

func LoadRegionById(db *sql.DB, id int) (*RegionLite, error) {
	trace(traceName(fmt.Sprintf("LoadRegionById(%d)", id)))
	rows, err := db.Query(
		"SELECT region_id, region_code, region_name, country_code, country_name, has_region_icon"+
			" FROM region_view"+
			" WHERE region_id=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, fmt.Errorf("Region not found (id: %d)", id)
	}
	region, err := readRegionFromRows(rows)
	if err != nil {
		return nil, err
	}
	return region, nil
}

func LoadRegionsByCountryCode(db *sql.DB, countryCode string) ([]RegionLite, error) {
	trace(traceName(fmt.Sprintf("LoadRegionsByCountryCode(%s)", countryCode)))
	rows, err := db.Query(
		"SELECT region_id, region_code, region_name, country_code, country_name, has_region_icon"+
			" FROM region_view"+
			" WHERE country_code=?"+
			" ORDER BY country_name, region_name", countryCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return readRegionListFromRows(rows)
}

func LoadRegionList(db *sql.DB) ([]RegionLite, error) {
	trace(traceName("LoadRegionList"))
	rows, err := db.Query(
		"SELECT region_id, region_code, region_name, country_code, country_name, has_region_icon" +
			" FROM region_view" +
			" ORDER BY country_name, region_name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return readRegionListFromRows(rows)
}

func DeleteRegion(db *sql.DB, regionId int) error {
	defer trace(traceName(fmt.Sprintf("DeleteRegion(%d)", regionId)))
	_, err := db.Exec("DELETE FROM regions WHERE region_id=?", regionId)
	return err
}

func InsertRegion(db *sql.DB, data RegionData) (*RegionLite, error) {
	defer trace(traceName(fmt.Sprintf("InsertRegion (data: %v)", data)))
	res, err := db.Exec(
		"INSERT INTO regions"+
			" (name, code, country_code)"+
			" VALUES(?, ?, ?)",
		data.Name, data.Code, data.CountryCode)
	if err != nil {
		return nil, err
	}
	regionId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return LoadRegionById(db, int(regionId))
}

func UpdateRegion(db *sql.DB, regionId int, data RegionData) error {
	defer trace(traceName(fmt.Sprintf("UpdateRegion(%d, %v)", regionId, data)))
	_, err := db.Exec(
		"UPDATE regions"+
			" SET name=?, code=?, country_code=?"+
			" WHERE id=?",
		data.Name, data.Code, data.CountryCode, regionId)
	return err
}

func readRegionListFromRows(rows *sql.Rows) ([]RegionLite, error) {
	var regions []RegionLite
	for rows.Next() {
		region, err := readRegionFromRows(rows)
		if err != nil {
			return nil, err
		}
		regions = append(regions, *region)
	}
	return regions, nil
}

func readRegionFromRows(rows *sql.Rows) (*RegionLite, error) {
	region := RegionLite{}
	err := rows.Scan(&region.Id, &region.Code, &region.Name, &region.CountryCode, &region.CountryName, &region.HasRegionIcon)
	if err != nil {
		return nil, err
	}
	return &region, nil
}
