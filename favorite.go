package main

import (
	"database/sql"
	"fmt"
)

func InsertFavorite(db *sql.DB, userId int, personId int) error {
	defer trace(traceName(fmt.Sprintf("InsertFavorite(%d, %d)", userId, personId)))
	_, err := db.Exec(
		"INSERT IGNORE INTO favorite_people"+
			" (user_id, person_id)"+
			" VALUES (?, ?)",
		userId, personId)
	return err
}

func DeleteFavorite(db *sql.DB, userId int, personId int) error {
	defer trace(traceName(fmt.Sprintf("DeleteFavorite(%d, %d)", userId, personId)))
	_, err := db.Exec(
		"DELETE FROM favorite_people"+
			" WHERE user_id=? AND person_id=?",
		userId, personId)
	return err
}
