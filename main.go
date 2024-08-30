package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	appPath, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}
	DbFile := filepath.Join(appPath, DBFile)

	err = checkDB(DbFile)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open(dbDriver, DBFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := chi.NewRouter()

	server := http.FileServer(http.Dir(webDir))
	router.Handle("/*", http.StripPrefix("/", server))
	router.Route("/api", func(r chi.Router) {
		r.Get("/nextdate", nextDateHandler)
		r.Get("/tasks", getTasksHandler(db))

		r.Route("/task", func(rout chi.Router) {
			rout.Get("/", getTaskHandler(db))
			rout.Post("/", addTaskHandler(db))
			rout.Delete("/", deleteTaskHandler(db))
			rout.Put("/", updateTaskHandler(db))
			rout.Post("/done", completeTaskHandler(db))
		})
	})
	if err = http.ListenAndServe(":"+strconv.Itoa(ServPort), router); err != nil {
		fmt.Printf("Error start server: %s", err.Error())
		return
	}
}
