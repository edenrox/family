package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func getIntPathParam(r *http.Request, paramName string, index int) (int, error) {
	param, err := getPathParam(r, paramName, index)
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(param)
}

func getPathParam(r *http.Request, paramName string, index int) (string, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) <= index {
		return "", fmt.Errorf("Path %s missing expected %s at index %d", r.URL.Path, paramName, index)
	}
	return parts[index], nil
}
