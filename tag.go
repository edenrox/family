package main

import (
	"database/sql"
	"fmt"
)

type Tag struct {
	Id    int
	Label string
}

func InsertTag(db *sql.DB, label string) (*Tag, error) {
	defer trace(traceName(fmt.Sprintf("InsertTag(%s)", label)))
	res, err := db.Exec(
		"INSERT INTO tags"+
			" (label)"+
			" VALUES (?)",
		label)
	tagId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	} else {
		return &Tag{Id: int(tagId), Label: label}, nil
	}
}

func LoadTagById(db *sql.DB, tagId int) (*Tag, error) {
	defer trace(traceName(fmt.Sprintf("LoadTagById(%d)", tagId)))
	rows, err := db.Query(
		"SELECT id, label"+
			" FROM tags"+
			" WHERE id=?",
		tagId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, fmt.Errorf("Tag not found (%+v)", tagId)
	}
	return readTagFromRows(rows)
}

func LoadTagByLabel(db *sql.DB, label string) (*Tag, error) {
	defer trace(traceName(fmt.Sprintf("LoadTagByLabel(%s)", label)))
	rows, err := db.Query(
		"SELECT id, label"+
			" FROM tags"+
			" WHERE label=?",
		label)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, fmt.Errorf("Tag not found (%+v)", label)
	}
	return readTagFromRows(rows)
}

func InsertPeopleTag(db *sql.DB, tagId int, personId int) error {
	defer trace(traceName(fmt.Sprintf("InsertPeopleTag(%d, %d)", tagId, personId)))
	_, err := db.Exec(
		"INSERT IGNORE INTO people_tags"+
			" (person_id, tag_id)"+
			" VALUES (?, ?)",
		personId, tagId)
	return err
}

func LoadTagsListByPrefix(db *sql.DB, labelPrefix string) ([]Tag, error) {
	defer trace(traceName("LoadTags()"))
	rows, err := db.Query(
		"SELECT id, label"+
			" FROM tags"+
			" WHERE label LIKE CONCAT(?, '%')"+
			" ORDER BY label",
		labelPrefix)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return readTagListFromRows(rows)
}

func LoadTagsForPerson(db *sql.DB, personId int) ([]Tag, error) {
	defer trace(traceName(fmt.Sprintf("LoadTagsForPerson(%d)", personId)))
	rows, err := db.Query(
		"SELECT t.id, t.label"+
			" FROM tags t"+
			"   INNER JOIN people_tags pt ON pt.tag_id=t.id"+
			" WHERE pt.person_id=?"+
			" ORDER BY t.label",
		personId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return readTagListFromRows(rows)
}

func DeletePeopleTag(db *sql.DB, tagId int, personId int) error {
	defer trace(traceName(fmt.Sprintf("DeletePeopleTag(%d, %d)", tagId, personId)))
	_, err := db.Exec(
		"DELETE FROM people_tags"+
			" WHERE tag_id=? AND person_id=?",
		tagId, personId)
	return err
}

func DeleteTag(db *sql.DB, tagId int) error {
	defer trace(traceName(fmt.Sprintf("DeleteTag(%d)", tagId)))
	tx, err := db.Begin()
	_, err = db.Exec(
		"DELETE FROM tags"+
			" WHERE id=?",
		tagId)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = db.Exec(
		"DELETE FROM people_tags"+
			" WHERE tag_id=?",
		tagId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func readTagListFromRows(rows *sql.Rows) ([]Tag, error) {
	var tags []Tag
	for rows.Next() {
		tag, err := readTagFromRows(rows)
		if err != nil {
			return nil, err
		}
		tags = append(tags, *tag)
	}
	return tags, nil
}

func readTagFromRows(rows *sql.Rows) (*Tag, error) {
	tag := Tag{}
	err := rows.Scan(&tag.Id, &tag.Label)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}
