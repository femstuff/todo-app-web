package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	appPath, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}
	DbFile := filepath.Join(appPath, DBFile)

	err = CheckDB(DbFile)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open(DbDriver, DBFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := NewSchedulerStore(db)

	router := chi.NewRouter()

	server := http.FileServer(http.Dir(webDir))
	router.Handle("/*", http.StripPrefix("/", server))
	router.Route("/api", func(r chi.Router) {
		r.Get("/nextdate", nextDateHandler)
		r.Get("/tasks", getTasksHandler(&store))

		r.Route("/task", func(rout chi.Router) {
			rout.Get("/", getTaskHandler(&store))
			rout.Post("/", addTaskHandler(&store))
			rout.Delete("/", deleteTaskHandler(&store))
			rout.Put("/", updateTaskHandler(&store))
			rout.Post("/done", completeTaskHandler(&store))
		})
	})

	log.Printf("Server start with port: %v \n", ServPort)
	if err = http.ListenAndServe(":"+strconv.Itoa(ServPort), router); err != nil {
		fmt.Printf("Error start server: %s", err.Error())
		return
	}
}
