package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/vrazinsky/go-final-project/calc"
	"github.com/vrazinsky/go-final-project/models"
)

func (h *Handlers) HandleAddTask(res http.ResponseWriter, req *http.Request) {
	var input models.AddTaskInput
	var buf bytes.Buffer
	var date time.Time
	var dateToDb string
	now := time.Now()
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		res.Write(ErrorAddTaskResponse(err, ""))
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &input); err != nil {
		res.Write(ErrorAddTaskResponse(err, ""))
		return
	}

	if input.Title == "" {
		res.Write(ErrorAddTaskResponse(nil, "no title"))
		return
	}

	if len(*input.Date) != 8 && len(*input.Date) != 0 {
		res.Write(ErrorAddTaskResponse(nil, "incorrect date format"))
		return
	}

	if input.Date == nil || *input.Date == "" {
		date = now
		dateToDb = now.Format("20060102")
	} else {
		date, err = time.Parse("20060102", *input.Date)
		if err != nil {
			res.Write(ErrorAddTaskResponse(nil, "incorrect date format"))
			return
		}
		dateToDb = *input.Date
	}
	if IsDateAfter(now, date) {
		if input.Repeat == nil || *input.Repeat == "" {
			dateToDb = now.Format("20060102")
		} else {
			dateToDb, err = calc.NextDate(now, date.Format("20060102"), *input.Repeat)
			if err != nil {
				res.Write(ErrorAddTaskResponse(err, ""))
				return
			}
		}
	}

	var insertId int64 = 0
	row := h.db.QueryRowContext(h.ctx, addTaskQuery,
		sql.Named("title", input.Title),
		sql.Named("date", dateToDb),
		sql.Named("comment", input.Comment),
		sql.Named("repeat", input.Repeat))
	err = row.Scan(&insertId)
	if err != nil {
		res.Write(ErrorAddTaskResponse(err, ""))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	response := models.AddTaskResponse{Id: insertId}
	responseBytes, _ := json.Marshal(response)
	res.Write(responseBytes)
}
