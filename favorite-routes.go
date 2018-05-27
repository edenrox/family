package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func favoriteAdd(w http.ResponseWriter, r *http.Request) {
	userId := 1
	personId, err := strconv.Atoi(r.FormValue("person_id"))
	if personId < 1 {
		http.Error(w, fmt.Sprintf("Error parsing person_id: %v", err), http.StatusBadRequest)
		return
	}

	err = InsertFavorite(db, userId, personId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting favorite: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func favoriteDelete(w http.ResponseWriter, r *http.Request) {
	userId := 1
	personId, err := strconv.Atoi(r.FormValue("person_id"))
	if personId < 1 {
		http.Error(w, fmt.Sprintf("Error parsing person_id: %v", err), http.StatusBadRequest)
		return
	}

	err = DeleteFavorite(db, userId, personId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting favorite: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func init() {
	http.HandleFunc("/favorite/add", favoriteAdd)
	http.HandleFunc("/favorite/delete", favoriteDelete)
}
