package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/vrazinsky/go-final-project/models"
	"github.com/vrazinsky/go-final-project/nextdate"
)

func (h *Handlers) HandleAddTask(res http.ResponseWriter, req *http.Request) {
	var input models.TaskInput
	var buf bytes.Buffer
	var date time.Time
	var dateToDb string
	now := time.Now()
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		logWriteErr(res.Write(ErrorResponse(err, "")))
	}
	if err = json.Unmarshal(buf.Bytes(), &input); err != nil {
		logWriteErr(res.Write(ErrorResponse(err, "")))
		return
	}

	if input.Title == "" {
		logWriteErr(res.Write(ErrorResponse(nil, "no title")))
		return
	}

	if len(*input.Date) != 8 && len(*input.Date) != 0 {
		logWriteErr(res.Write(ErrorResponse(nil, "incorrect date format")))
		return
	}

	if input.Date == nil || *input.Date == "" {
		date = now
		dateToDb = now.Format(layout)
	} else {
		date, err = time.Parse("20060102", *input.Date)
		if err != nil {
			logWriteErr(res.Write(ErrorResponse(nil, "incorrect date format")))
			return
		}
		dateToDb = *input.Date
	}
	if IsDateAfter(now, date) {
		if input.Repeat == nil || *input.Repeat == "" {
			dateToDb = now.Format(layout)
		} else {
			dateToDb, err = nextdate.NextDate(now, date.Format(layout), *input.Repeat)
			if err != nil {
				logWriteErr(res.Write(ErrorResponse(err, "")))
				return
			}
		}
	}
	var repeat, comment string
	if input.Repeat == nil {
		repeat = ""
	} else {
		repeat = *input.Repeat
	}

	if input.Comment == nil {
		comment = ""
	} else {
		comment = *input.Comment
	}

	task := models.Task{
		Date:    dateToDb,
		Title:   input.Title,
		Comment: comment,
		Repeat:  repeat,
	}

	insertId, err := h.db.AddTask(task)
	if err != nil {
		logWriteErr(res.Write(ErrorResponse(err, "")))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	response := models.Task{Id: insertId}
	responseBytes, _ := json.Marshal(response)
	logWriteErr(res.Write(responseBytes))
}
