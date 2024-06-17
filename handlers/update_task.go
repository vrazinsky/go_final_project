package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/vrazinsky/go-final-project/models"
	"github.com/vrazinsky/go-final-project/nextdate"
)

func (h *Handlers) HandleUpdateTask(res http.ResponseWriter, req *http.Request) {
	var task models.UpdateTaskInput
	var buf bytes.Buffer
	var date time.Time
	var dateToDb string
	now := time.Now()
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		logWriteErr(res.Write(ErrorResponse(err, "")))
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		logWriteErr(res.Write(ErrorResponse(err, "")))
		return
	}

	if len(task.Id) == 0 {
		logWriteErr(res.Write(ErrorAddTaskResponse(nil, "no id")))
		return
	}
	_, err = strconv.Atoi(task.Id)
	if err != nil {
		logWriteErr(res.Write(ErrorAddTaskResponse(nil, "incorrect id")))
		return
	}
	if task.Title == "" {
		logWriteErr(res.Write(ErrorAddTaskResponse(nil, "no title")))
		return
	}

	if len(*task.Date) != 8 && len(*task.Date) != 0 {
		logWriteErr(res.Write(ErrorAddTaskResponse(nil, "incorrect date format")))
		return
	}

	if task.Date == nil || *task.Date == "" {
		date = now
		dateToDb = now.Format("20060102")
	} else {
		date, err = time.Parse("20060102", *task.Date)
		if err != nil {
			logWriteErr(res.Write(ErrorAddTaskResponse(nil, "incorrect date format")))
			return
		}
		dateToDb = *task.Date
	}
	if IsDateAfter(now, date) {
		if task.Repeat == nil || *task.Repeat == "" {
			dateToDb = now.Format("20060102")
		} else {
			dateToDb, err = nextdate.NextDate(now, date.Format("20060102"), *task.Repeat)
			if err != nil {
				logWriteErr(res.Write(ErrorAddTaskResponse(err, "")))
				return
			}
		}
	}

	result, err := h.db.ExecContext(h.ctx, updateTakQuery,
		sql.Named("id", task.Id),
		sql.Named("title", task.Title),
		sql.Named("date", dateToDb),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		logWriteErr(res.Write(ErrorAddTaskResponse(err, "")))
		return
	}
	updatedRowsNumber, err := result.RowsAffected()
	if err != nil {
		logWriteErr(res.Write(ErrorAddTaskResponse(err, "")))
		return
	}
	if updatedRowsNumber == 0 {
		logWriteErr(res.Write(ErrorAddTaskResponse(nil, "task not found")))
		return
	}
	var response struct{}
	responseBytes, _ := json.Marshal(response)
	logWriteErr(res.Write(responseBytes))
}
