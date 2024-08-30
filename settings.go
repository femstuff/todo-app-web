package main

const (
	dbFile      = "./scheduler.db"
	dbDriver    = "sqlite3"
	dateFormat  = "20060102"
	webDir      = "./web/"
	defaultPort = 7540
)

var ServPort = getServPort("TODO_PORT")
var DBFile = getDbFile("TODO_DBFILE")
