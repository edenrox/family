package main

import (
	"database/sql"
	"fmt"
	"log"
)

type RegionLite struct {
	Id          int
	Code        string
	Name        string
	CountryCode string
	CountryName string
}

func LoadRegionById(db *sql.DB, id int) (*RegionLite, error) {
	log.Printf("Load RegionLite id: %d", id)
	rows, err := db.Query("SELECT region_id, region_code, region_name, country_code, country_name FROM region_view WHERE region_id=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, fmt.Errorf("Region not found (id: %d)", id)
	}
	region := RegionLite{}
	err = rows.Scan(&region.Id, &region.Code, &region.Name, &region.CountryCode, &region.CountryName)
	if err != nil {
		return nil, err
	}
	return &region, nil
}

func LoadRegionList(db *sql.DB) ([]RegionLite, error) {
	log.Printf("Load RegionLite list")
	rows, err := db.Query("SELECT region_id, region_code, region_name, country_code, country_name FROM region_view")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var regions []RegionLite
	for rows.Next() {
		region := RegionLite{}
		err = rows.Scan(&region.Id, &region.Code, &region.Name, &region.CountryCode, &region.CountryName)
		if err != nil {
			return nil, err
		}
		regions = append(regions, region)
	}
	return regions, nil
}
