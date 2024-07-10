package task

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

type TaskStorage struct {
	db *sqlx.DB
}

func NewTaskStorage() (TaskStorage, error) {
	err := CreateDB()
	if err != nil {
		return TaskStorage{}, err
	}

	db, err := sqlx.Connect("sqlite3", filepath.Join(".", "scheduler.db"))
	if err != nil {
		return TaskStorage{}, err
	}

	return TaskStorage{db: db}, nil
}

func CreateDB() error {
	dbFileName := "scheduler.db"
	dbFile := filepath.Join(".", dbFileName)

	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath != "" {
		dbFile = dbPath
		_, err := os.Stat(dbFile)
		if err == nil {
			fmt.Println("Database already exists at TODO_DBFILE:", dbFile)
			return nil
		}
	}

	_, err := os.Stat(dbFile)
	switch {
	case err != nil:
		err := createDatabase(dbFile)
		if err != nil {
			return err
		}
		if dbPath != "" {
			fmt.Println("Database created successfully at TODO_DBFILE:", dbFile)
		} else {
			fmt.Println("Database created successfully at:", dbFile)
		}
	default:
		fmt.Println("Database already exists at:", dbFile)
		return nil
	}

	return nil
}

func createDatabase(dbFile string) error {
	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(128) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT ""
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);`
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		return err
	}

	return nil
}
