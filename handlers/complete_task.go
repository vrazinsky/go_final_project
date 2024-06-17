package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/vrazinsky/go-final-project/models"
	"github.com/vrazinsky/go-final-project/nextdate"
)

func (h *Handlers) HandleCompleteTask(res http.ResponseWriter, req *http.Request) {
	idStr := req.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logWriteErr(res.Write(ErrorResponse(nil, "incorrect input data")))
		return
	}
	task, err := h.db.GetTask(id)
	if err != nil {
		logWriteErr(res.Write(ErrorResponse(err, "")))
		return
	}
	if task.Repeat == "" {
		err = h.db.DeleteTask(id)
		if err != nil {
			logWriteErr(res.Write(ErrorResponse(err, "")))
			return
		}
		var response struct{}
		responseBytes, _ := json.Marshal(response)
		logWriteErr(res.Write(responseBytes))
		return
	}
	nextDate, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		logWriteErr(res.Write(ErrorResponse(err, "")))
		return
	}

	nextTask := models.Task{
		Id:      task.Id,
		Date:    nextDate,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	err = h.db.UpdateTask(nextTask)
	if err != nil {
		logWriteErr(res.Write(ErrorResponse(err, "")))
		return
	}
	var response struct{}
	responseBytes, _ := json.Marshal(response)
	logWriteErr(res.Write(responseBytes))
}
