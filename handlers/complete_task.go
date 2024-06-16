package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/vrazinsky/go-final-project/calc"
	"github.com/vrazinsky/go-final-project/models"
)

func (h *Handlers) HandleCompleteTask(res http.ResponseWriter, req *http.Request) {
	idStr := req.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res.Write(ErrorResponse(nil, "incorrect input data"))
		return
	}
	row := h.db.QueryRowContext(h.ctx, getTaskQuery, sql.Named("id", id))
	var task models.Task
	var taskId int
	err = row.Scan(&taskId, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	task.Id = strconv.Itoa(taskId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			res.Write(ErrorResponse(nil, "task not found"))
		} else {
			res.Write(ErrorResponse(err, ""))
		}
		return
	}
	if task.Repeat == nil || *task.Repeat == "" {
		_, err = h.db.ExecContext(h.ctx, deleteTaskQuery, sql.Named("id", id))
		if err != nil {
			res.Write(ErrorResponse(err, ""))
			return
		}
		var response struct{}
		responseBytes, _ := json.Marshal(response)
		res.Write(responseBytes)
		return
	}
	nextDate, err := calc.NextDate(time.Now(), task.Date, *task.Repeat)
	if err != nil {
		res.Write(ErrorResponse(err, ""))
		return
	}
	_, err = h.db.ExecContext(h.ctx, updateTakQuery,
		sql.Named("id", task.Id),
		sql.Named("title", task.Title),
		sql.Named("date", nextDate),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		res.Write(ErrorAddTaskResponse(err, ""))
		return
	}
	var response struct{}
	responseBytes, _ := json.Marshal(response)
	res.Write(responseBytes)
}
