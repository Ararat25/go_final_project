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

type Task struct {
	ID      int64  `json:"id" db:"id"`
	Date    string `json:"date" db:"date"`
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}

// SchedulerStore структура для работы с базой данных scheduler
type SchedulerStore struct {
	db *sql.DB
}

// Connect присоединяется к базе данных и возвращает *SchedulerStore
func Connect() (*SchedulerStore, error) {
	exist := true
	if checkExistsDBFile() {
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
		err = store.createSchedulerTable()
		if err != nil {
			return nil, err
		}

		err = store.createIndexDate()
		if err != nil {
			return nil, err
		}
	}
	return &store, nil
}

// Close закрывает соединение с бд
func (db *SchedulerStore) Close() {
	db.db.Close()
}

// checkExistsDBFile проверяет существует ли файл с бд, если true значит файла не существует
func checkExistsDBFile() bool {
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

// createSchedulerTable создает scheduler таблицу в базе данных
func (db *SchedulerStore) createSchedulerTable() error {
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

// CreateSchedulerTable создает scheduler таблицу в базе данных
func (db *SchedulerStore) createIndexDate() error {
	sqlCreateIndexQuery := `CREATE INDEX IF NOT EXISTS date_index ON scheduler (date);`

	_, err := db.ExecuteQuery(sqlCreateIndexQuery)
	if err != nil {
		return err
	}

	return nil
}

// AddTask добавляет задачу в базу данных
func (db *SchedulerStore) AddTask(task Task) (int64, error) {
	result, err := db.db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) 
	VALUES (:date, :title, :comment, :repeat)`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}
