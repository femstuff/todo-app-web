package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func zeroTime(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func checkDB(dbPath string) error {
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		db, err := sql.Open(dbDriver, dbPath)
		if err != nil {
			return err
		}
		defer db.Close()

		query := `
        CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date VARCHAR(8) NOT NULL,
			title TEXT NOT NULL,
			comment TEXT DEFAULT "",
			repeat VARCHAR(128) DEFAULT ""
		);
		CREATE INDEX idx_scheduler_date ON scheduler(date);
        `
		_, err = db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}

		return nil
	} else if err != nil {
		return err
	}

	return nil
}

func NextDate(now time.Time, date, repeat string) (string, error) {
	parsDate, err := time.Parse(dateFormat, date)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %w", err)
	}

	var nextDate time.Time
	var days int
	check := strings.Split(repeat, " ")

	if len(check) > 2 {
		return "", fmt.Errorf("неверный формат правил повторения")
	}

	switch check[0] {
	case "y":
		nextDate = parsDate.AddDate(1, 0, 0)
	case "d":
		if len(check) != 2 {
			return "", fmt.Errorf("не указан второй аргумент для дней")
		}

		days, err = strconv.Atoi(check[1])
		if err != nil {
			return "", fmt.Errorf("не удалось преобразовать количество дней: %w", err)
		}

		if days >= 1 && days <= 400 {
			nextDate = parsDate.AddDate(0, 0, days)
		} else {
			return "", fmt.Errorf("неверный диапазон заданных дней: %d", days)
		}
	case "":
		return "", fmt.Errorf("правила повторения не указаны")
	default:
		return "", fmt.Errorf("некорректное правило повторения: %s", check[0])

	}

	for nextDate.Before(now) {
		if check[0] == "d" {
			nextDate = nextDate.AddDate(0, 0, days)
		} else if check[0] == "y" {
			nextDate = nextDate.AddDate(1, 0, 0)
		} else {
			return "", fmt.Errorf("некорректное правило повторения: %s", check[0])
		}
	}

	return nextDate.Format(dateFormat), nil
}

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
