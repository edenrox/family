package main

import (
	"log"
	"time"
)

func traceName(name string) (string, time.Time) {
	log.Printf("%s - start", name)
	return name, time.Now()
}

func trace(name string, start time.Time) {
	elapsed := time.Since(start)
	log.Printf("%s - end took: %s", name, elapsed)
}
