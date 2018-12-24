package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func tagJsonAdd(w http.ResponseWriter, r *http.Request) {
	label := strings.TrimSpace(r.FormValue("label"))
	if label == "" {
		http.Error(w, fmt.Sprintf("Error label must be non-empty: %v", label), http.StatusBadRequest)
		return
	}
	personId, err := strconv.Atoi(r.FormValue("person_id"))
	if personId < 1 {
		http.Error(w, fmt.Sprintf("Error parsing person_id: %v", err), http.StatusBadRequest)
		return
	}
	tag, err := LoadTagByLabel(db, label)
	if tag == nil {
		tag, err = InsertTag(db, label)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error inserting tag: %v", err), http.StatusInternalServerError)
			return
		}
	}

	err = InsertPeopleTag(db, tag.Id, personId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting tag: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func tagJsonDelete(w http.ResponseWriter, r *http.Request) {
	label := strings.TrimSpace(r.FormValue("label"))
	tagId, err := strconv.Atoi(r.FormValue("tag_id"))
	if label == "" && tagId == 0 {
		http.Error(w, fmt.Sprintf("Supply either label (%v) or tagId (%v)", label, tagId), http.StatusBadRequest)
		return
	}
	personId, err := strconv.Atoi(r.FormValue("person_id"))
	if personId < 1 {
		http.Error(w, fmt.Sprintf("Error parsing person_id: %v", err), http.StatusBadRequest)
		return
	}
	if tagId == 0 {
		tag, _ := LoadTagByLabel(db, label)
		if tag == nil {
			http.Error(w, fmt.Sprintf("Tag with specified label not found: %v", label), http.StatusBadRequest)
			return
		}
		tagId = tag.Id
	}
	err = DeletePeopleTag(db, tagId, personId)
	w.WriteHeader(http.StatusAccepted)
}

func tagJsonList(w http.ResponseWriter, r *http.Request) {
	labelPrefix := strings.TrimSpace(r.FormValue("prefix"))
	tags, err := LoadTagsListByPrefix(db, labelPrefix)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading tags: %v", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

func tagList(w http.ResponseWriter, r *http.Request) {
	data, err := LoadTagsListByPrefix(db, "")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading tags: %v", err), 500)
		return
	}
	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/tag/list.html")).Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func tagView(w http.ResponseWriter, r *http.Request) {
	tagLabel, err := getPathParam(r, "tagLabel", 3)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing tagId: %v", err), 400)
		return
	}

	tag, err := LoadTagByLabel(db, tagLabel)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading tag: %v", err), 400)
		return
	}
	personList, err := LoadPersonLiteListWithTag(db, tag.Label)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading person list: %v", err), 400)
		return
	}

	data := struct {
		Tag    *Tag
		People []PersonLite
	}{
		tag,
		personList,
	}

	err = template.Must(template.ParseFiles("tmpl/layout/main.html", "tmpl/tag/view.html")).Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func addTagRoutes() {
	http.HandleFunc("/tag/json/add", tagJsonAdd)
	http.HandleFunc("/tag/json/delete", tagJsonDelete)
	http.HandleFunc("/tag/json/list", tagJsonList)
	http.HandleFunc("/tag/list", tagList)
	http.HandleFunc("/tag/view/", tagView)
}
