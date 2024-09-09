package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

type SchedulerStore struct {
	db *sql.DB
}

func NewSchedulerStore(db *sql.DB) SchedulerStore {
	return SchedulerStore{db: db}
}

func (s SchedulerStore) Add(t Task) (int, error) {
	res, err := s.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))

	if err != nil {
		return 0, err
	}

	lastInserted, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastInserted), nil
}

func (s SchedulerStore) Update(t Task) error {
	_, err := s.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
		sql.Named("id", t.ID))

	if err != nil {
		return err
	}

	return nil
}

func (s SchedulerStore) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	if err != nil {
		return err
	}

	return nil
}

func (s SchedulerStore) Get(id string) (Task, error) {
	var task Task

	row := s.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func (s SchedulerStore) GetTasks() (TaskResponse, error) {
	var tasks TaskResponse

	rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit",
		sql.Named("limit", dbLimit))
	if err != nil {
		return TaskResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return TaskResponse{}, err
		}

		tasks.Tasks = append(tasks.Tasks, task)
	}
	if err = rows.Err(); err != nil {
		return TaskResponse{}, err
	}

	if len(tasks.Tasks) == 0 {
		tasks.Tasks = []Task{}
		return tasks, nil
	}

	return tasks, nil
}
