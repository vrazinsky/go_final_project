package models

import (
	"errors"
	"time"

	"github.com/vrazinsky/go-final-project/internal/nextdate"
	"github.com/vrazinsky/go-final-project/internal/utils"
)

type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type TaskInput struct {
	Id      string  `json:"id"`
	Date    *string `json:"date,omitempty"`
	Title   string  `json:"title"`
	Comment *string `json:"comment,omitempty"`
	Repeat  *string `json:"repeat,omitempty"`
}

func (i *TaskInput) Validate() error {
	if i.Title == "" {
		return errors.New("no title")
	}
	if i.Repeat != nil && *i.Repeat != "" {
		_, _, _, err := nextdate.ParseRepeat(*i.Repeat)
		if err != nil {
			return err
		}
	}
	if i.Date != nil && *i.Date != "" {
		_, err := time.Parse(utils.Layout, *i.Date)
		if err != nil {
			return err
		}
	}
	return nil
}

type GetTasksResponse struct {
	Error string  `json:"error,omitempty"`
	Tasks *[]Task `json:"tasks"`
}

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

type SignInInput struct {
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}
