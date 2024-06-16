package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/vrazinsky/go-final-project/models"
)

func (h *Handlers) HandleGetTask(res http.ResponseWriter, req *http.Request) {
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
	data, _ := json.Marshal(task)
	res.Write(data)
}
