package main

import (
	"log"
	"os"
	"strconv"
)

const (
	dbFile      = "./scheduler.db"
	DbDriver    = "sqlite3"
	dateFormat  = "20060102"
	webDir      = "./web/"
	defaultPort = 7540
	dbLimit     = 50
)

var ServPort = getServPort("TODO_PORT")
var DBFile = getDbFile("TODO_DBFILE")

func getServPort(env string) int {
	key := os.Getenv(env)
	if key == "" {
		return defaultPort
	}
	port, err := strconv.Atoi(key)
	if err != nil {
		log.Fatal(err.Error())
	}
	return port
}

func getDbFile(env string) string {
	key := os.Getenv(env)

	if key == "" {
		return dbFile
	}
	return key
}
