package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/vrazinsky/go-final-project/models"
)

type Handlers struct {
	db  *sql.DB
	ctx context.Context
}

func NewHandler(ctx context.Context, db *sql.DB) *Handlers {
	return &Handlers{ctx: ctx, db: db}
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

func logWriteErr(_ int, err error) {
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}
