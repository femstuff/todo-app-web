package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type TaskResponse struct {
	Tasks []Task `json:"tasks"`
}

func JsonError(w http.ResponseWriter, messageError string, codeError int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codeError)
	errRes := Response{
		Error: messageError,
	}
	json.NewEncoder(w).Encode(errRes)
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if len(now) == 0 || len(date) == 0 || len(repeat) == 0 {
		http.Error(w, "empty now, date or repeat", http.StatusBadRequest)
		return
	}

	nowPars, err := time.Parse(dateFormat, now)
	if err != nil {
		fmt.Errorf("Error with parsing: %s", err.Error())
		return
	}
	nextDate, err := NextDate(nowPars, date, repeat)
	if err != nil {
		http.Error(w, "err with request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}

func getTaskHandler(store *SchedulerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			JsonError(w, `error`, http.StatusBadRequest)
			return
		}

		task, err := store.Get(id)
		if err != nil {
			JsonError(w, `error: Task not found`, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(task)
	}
}

func deleteTaskHandler(store *SchedulerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			JsonError(w, `error`, http.StatusBadRequest)
			return
		}

		task, err := store.Get(id)
		if err != nil {
			JsonError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if task.Repeat == "" {
			err = store.Delete(task.ID)
			if err != nil {
				JsonError(w, `error with delete task`, http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{}`))
			return
		}

		err = store.Delete(task.ID)
		if err != nil {
			JsonError(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}
}

func addTaskHandler(store *SchedulerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			JsonError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			JsonError(w, "error with deserialize json", http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			JsonError(w, "title couldnot be empty", http.StatusBadRequest)
			return
		}

		now := time.Now()

		if task.Date == "" {
			task.Date = now.Format(dateFormat)
		}

		rep := strings.Split(task.Repeat, " ")

		if len(rep[0]) != 0 && (rep[0] != "d" && rep[0] != "y") {
			JsonError(w, "incorrect rules repeat", http.StatusBadRequest)
			return
		}

		date, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			JsonError(w, "error with parsing date", http.StatusBadRequest)
			return
		}

		nowZero := zeroTime(now)
		if date.Before(nowZero) {
			if task.Repeat == "" {
				task.Date = now.Format(dateFormat)
			} else {
				nextDate, err := NextDate(now, date.Format(dateFormat), task.Repeat)
				if err != nil {
					JsonError(w, err.Error(), http.StatusBadRequest)
					return
				}

				task.Date = nextDate
			}
		}

		res, err := store.Add(task)
		if err != nil {
			JsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resulTask, err := store.Get(strconv.Itoa(res))
		if err != nil {
			JsonError(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resulTask)
	}
}

func updateTaskHandler(store *SchedulerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			JsonError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			JsonError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if task.ID == "" || task.Date == "" || task.Title == "" {
			JsonError(w, `task not found`, http.StatusBadRequest)
			return
		}

		now := time.Now()

		date, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			JsonError(w, `incorrect date`, http.StatusBadRequest)
			return
		}

		if date.Before(now) {
			if task.Repeat == "" {
				task.Date = now.Format(dateFormat)
			} else {
				nextDate, err := NextDate(now, date.Format(dateFormat), task.Repeat)

				if err != nil {
					JsonError(w, err.Error(), http.StatusBadRequest)
					return
				}

				task.Date = nextDate
			}
		}

		_, err = store.Get(task.ID)
		if err != nil {
			JsonError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = store.Update(task)
		if err != nil {
			JsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}
}

func getTasksHandler(store *SchedulerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := store.GetTasks()
		if err != nil {
			JsonError(w, `error with get tasks`, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func completeTaskHandler(store *SchedulerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			JsonError(w, `error: empty id`, http.StatusBadRequest)
			return
		}

		task, err := store.Get(id)
		if err != nil {
			JsonError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if task.Repeat == "" {
			err = store.Delete(task.ID)
			if err != nil {
				JsonError(w, `error with delete task`, http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{}`))
			return
		}

		now := time.Now()

		taskDate, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			JsonError(w, `error: incorrect date`, http.StatusBadRequest)
			return
		}

		nextDate, err := NextDate(now, taskDate.Format(dateFormat), task.Repeat)
		if err != nil {
			JsonError(w, err.Error(), http.StatusBadRequest)
			return
		}
		task.Date = nextDate

		err = store.Update(task)
		if err != nil {
			JsonError(w, `error: cant update task`, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}
}
