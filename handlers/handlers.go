package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/vrazinsky/go-final-project/calc"
	"github.com/vrazinsky/go-final-project/models"
)

type Handlers struct {
	db  *sql.DB
	ctx context.Context
}

func NewHandler(ctx context.Context, db *sql.DB) *Handlers {
	return &Handlers{ctx: ctx, db: db}
}

func (h *Handlers) HandleNextTime(res http.ResponseWriter, req *http.Request) {
	now := req.FormValue("now")
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")
	if now == "" || date == "" || repeat == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	nowDate, err := time.Parse("20060102", now)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	nextDate, err := calc.NextDate(nowDate, date, repeat)
	if err != nil {
		res.Write([]byte(err.Error()))
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Write([]byte(nextDate))
}

func ErrorAddTaskResponse(err error, errMsg string) []byte {
	var response models.AddTaskResponse
	if err != nil {
		response = models.AddTaskResponse{Error: err.Error()}
	} else {
		response = models.AddTaskResponse{Error: errMsg}
	}
	data, _ := json.Marshal(response)
	return data
}

func ErrorGetTasksResponse(err error, errMsg string) []byte {
	var response models.GetTasksResponse
	if err != nil {
		response = models.GetTasksResponse{Error: err.Error()}
	} else {
		response = models.GetTasksResponse{Error: errMsg}
	}
	data, _ := json.Marshal(response)
	return data
}

func ErrorResponse(err error, errMsg string) []byte {
	var response models.ErrorResponse
	if err != nil {
		response = models.ErrorResponse{Error: err.Error()}
	} else {
		response = models.ErrorResponse{Error: errMsg}
	}
	data, _ := json.Marshal(response)
	return data
}

func IsDateAfter(date1, date2 time.Time) bool {
	date1Int, _ := strconv.Atoi(date1.Format("20060102"))
	date2Int, _ := strconv.Atoi(date2.Format("20060102"))
	return date1Int > date2Int

}
