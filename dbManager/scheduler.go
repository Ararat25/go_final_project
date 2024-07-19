package dbManager

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"

	"github.com/Ararat25/go_final_project/tests"
)

var dbPath = tests.DBFile

type SchedulerStore struct {
	db *sql.DB
}

// Connect присоединяется к базе данных и возвращает *SchedulerStore
func Connect() (*SchedulerStore, error) {
	exist := true
	if CheckExistsDBFile() {
		exist = false
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	store := SchedulerStore{
		db: db,
	}

	if !exist {
		store.CreateSchedulerTable()
	}
	return &store, nil
}

func (db *SchedulerStore) Close() {
	db.db.Close()
}

// CheckExistsDBFile проверяет существует ли файл с бд, если true значит файла не существует
func CheckExistsDBFile() bool {
	envFile := os.Getenv("TODO_DBFILE")
	if len(envFile) > 0 {
		dbPath = envFile
	}

	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := filepath.Join(filepath.Dir(appPath), dbPath)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	return install
}

// ExecuteQuery выполняет запрос, переданный в аргументе, к бд
func (db *SchedulerStore) ExecuteQuery(query string) (sql.Result, error) {
	res, err := db.db.Exec(query)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateSchedulerTable создает scheduler таблицу в базе данных
func (db *SchedulerStore) CreateSchedulerTable() error {
	sqlCreateTableQuery := `CREATE TABLE IF NOT EXISTS scheduler (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"date" CHAR(8) NOT NULL DEFAULT "",
		"title" VARCHAR(128) NOT NULL DEFAULT "",
		"comment" VARCHAR(256) NOT NULL DEFAULT "",
		"repeat" VARCHAR(128) NOT NULL DEFAULT ""
	);`

	_, err := db.ExecuteQuery(sqlCreateTableQuery)
	if err != nil {
		return err
	}

	return nil
}
