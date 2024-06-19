package handlers

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/vrazinsky/go-final-project/internal/models"
	"github.com/vrazinsky/go-final-project/internal/store"
	"github.com/vrazinsky/go-final-project/internal/utils"
)

type Handlers struct {
	db  store.DbService
	ctx context.Context
}

func NewHandler(ctx context.Context, db store.DbService) *Handlers {
	return &Handlers{ctx: ctx, db: db}
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
	date1Int, _ := strconv.Atoi(date1.Format(utils.Layout))
	date2Int, _ := strconv.Atoi(date2.Format(utils.Layout))
	return date1Int > date2Int

}

func logWriteErr(_ int, err error) {
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}
