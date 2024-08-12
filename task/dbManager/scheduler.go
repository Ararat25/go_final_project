package dbManager

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Ararat25/go_final_project/errors"
	"github.com/Ararat25/go_final_project/model"
	"github.com/Ararat25/go_final_project/model/entity"
	_ "modernc.org/sqlite"
)

var dbPath = "../scheduler.db"

const limit = 30 // лимит для количества возвращаемых задач из бд

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

// createSchedulerTable создает scheduler таблицу в базе данных
func (db *SchedulerStore) createSchedulerTable() error {
	sqlCreateTableQuery := `CREATE TABLE IF NOT EXISTS scheduler (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"date" CHAR(8) NOT NULL DEFAULT "",
		"title" VARCHAR(128) NOT NULL DEFAULT "",
		"comment" VARCHAR(256) NOT NULL DEFAULT "",
		"repeat" VARCHAR(128) NOT NULL DEFAULT ""
	);`

	_, err := db.db.Exec(sqlCreateTableQuery)
	if err != nil {
		return err
	}

	return nil
}

// createIndexDate создает index таблицу в базе данных по дате
func (db *SchedulerStore) createIndexDate() error {
	sqlCreateIndexQuery := `CREATE INDEX IF NOT EXISTS date_index ON scheduler (date);`

	_, err := db.db.Exec(sqlCreateIndexQuery)
	if err != nil {
		return err
	}

	return nil
}

// Close закрывает соединение с бд
func (db *SchedulerStore) Close() {
	db.db.Close()
}

// AddTask добавляет задачу в базу данных
func (db *SchedulerStore) AddTask(task entity.Task) (int64, error) {
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

// Find возвращает отфильтрованные по заданной строке задачи
func (db *SchedulerStore) Find(filter string) ([]entity.Task, error) {
	if filter == "" {
		tasks, err := db.getTasks()
		if err != nil {
			return nil, err
		}

		return tasks, nil
	} else {
		date, err := time.Parse("02.01.2006", filter)
		if err == nil {
			tasks, err := db.getTasksByDate(date.Format(model.TimeLayout))
			if err != nil {
				return nil, err
			}

			return tasks, nil
		} else {
			tasks, err := db.getTasksBySearchString(filter)
			if err != nil {
				return nil, err
			}

			return tasks, nil
		}
	}
}

// GetTaskById возвращает задачу по id
func (db *SchedulerStore) GetTaskById(id int) (*entity.Task, error) {
	rows, err := db.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler 
                                        WHERE id = :id`,
		sql.Named("id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newTask entity.Task

	rows.Next()

	err = rows.Scan(&newTask.ID, &newTask.Date, &newTask.Title, &newTask.Comment, &newTask.Repeat)
	if err != nil {
		return nil, err
	}

	return &newTask, nil
}

// EditTaskById изменяет пареметры задачи
func (db *SchedulerStore) EditTaskById(task *entity.Task) error {
	result, err := db.db.Exec(`UPDATE scheduler SET 
                     date = :date, 
                     title = :title, 
                     comment = :comment, 
                     repeat = :repeat WHERE id = :id`,
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return err
	}

	numRowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if numRowAffected == 0 {
		return errors.ErrNotValidID
	}

	return nil
}

// DeleteTask удаляет задачу из бд
func (db *SchedulerStore) DeleteTask(id int) error {
	result, err := db.db.Exec(`DELETE FROM scheduler WHERE id = :id`,
		sql.Named("id", id))
	if err != nil {
		return err
	}

	numRowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if numRowAffected == 0 {
		return errors.ErrNotValidID
	}

	return nil
}

// getTasks возвращает все задачи из бд
func (db *SchedulerStore) getTasks() ([]entity.Task, error) {
	rows, err := db.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit`,
		sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks, err := getTasksFromRows(rows)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// getTasksByDate возвращает задачи из бд по дате
func (db *SchedulerStore) getTasksByDate(date string) ([]entity.Task, error) {
	rows, err := db.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler 
                                        WHERE date = :date
                                        ORDER BY date
                                        LIMIT :limit`,
		sql.Named("limit", limit),
		sql.Named("date", date))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks, err := getTasksFromRows(rows)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// getTasksBySearchString возвращает задачи из бд по введенной строке
func (db *SchedulerStore) getTasksBySearchString(search string) ([]entity.Task, error) {
	rows, err := db.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler 
                                        WHERE title LIKE :search 
                                        OR comment LIKE :search
                                        ORDER BY date 
                                        LIMIT :limit `,
		sql.Named("limit", limit),
		sql.Named("search", fmt.Sprintf("%%%s%%", search)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks, err := getTasksFromRows(rows)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// getTasksFromRows возвращает данные из sql ответа в виде слайса Task
func getTasksFromRows(rows *sql.Rows) ([]entity.Task, error) {
	var tasks []entity.Task

	for rows.Next() {
		var newTask entity.Task

		err := rows.Scan(&newTask.ID, &newTask.Date, &newTask.Title, &newTask.Comment, &newTask.Repeat)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, newTask)
	}

	err := rows.Err()
	if err != nil {
		return nil, err
	}

	return tasks, nil
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
