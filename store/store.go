package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/vrazinsky/go-final-project/models"
	_ "modernc.org/sqlite"
)

type DbService struct {
	dbfile string
	db     *sql.DB
	ctx    context.Context
}

func NewDbService(dbfile string, ctx context.Context) *DbService {
	return &DbService{dbfile: dbfile, ctx: ctx}
}

func (s *DbService) InitDb() error {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := filepath.Join(filepath.Dir(appPath), s.dbfile)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	db, err := sql.Open("sqlite", s.dbfile)
	if err != nil {
		log.Fatal(err)
	}
	if install {
		_, err = db.ExecContext(s.ctx, createTableQuery)
		if err != nil {
			return err
		}
		_, err = db.ExecContext(s.ctx, createIndexQuery)
		if err != nil {
			return err
		}
	}
	s.db = db
	return nil
}

func (s *DbService) AddTask(input models.Task) (string, error) {
	var insertId int64 = 0
	row := s.db.QueryRowContext(s.ctx, addTaskQuery,
		sql.Named("title", input.Title),
		sql.Named("date", input.Date),
		sql.Named("comment", input.Comment),
		sql.Named("repeat", input.Repeat))
	err := row.Scan(&insertId)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(insertId, 10), nil
}

func (s *DbService) CompleteTask(id int) {

}

func (s *DbService) GetTask(id int) (models.Task, error) {
	row := s.db.QueryRowContext(s.ctx, getTaskQuery, sql.Named("id", id))
	var task models.Task
	err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func (s *DbService) GetTasks(filterByTitle, filterByDate bool, searchValue string) ([]models.Task, error) {
	tasks := make([]models.Task, 0)
	rows, err := s.db.QueryContext(s.ctx, getTasksQuery,
		sql.Named("filterByTitle", filterByTitle),
		sql.Named("filterByDate", filterByDate),
		sql.Named("searchValue", searchValue))
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *DbService) DeleteTask(id int) error {
	_, err := s.db.ExecContext(s.ctx, deleteTaskQuery, sql.Named("id", id))
	if err != nil {
		return err
	}
	return nil
}

func (s *DbService) UpdateTask(task models.Task) error {
	result, err := s.db.ExecContext(s.ctx, updateTakQuery,
		sql.Named("id", task.Id),
		sql.Named("title", task.Title),
		sql.Named("date", task.Date),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return err
	}
	updatedRowsNumber, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if updatedRowsNumber == 0 {
		return errors.New("task not found")
	}
	return nil
}
