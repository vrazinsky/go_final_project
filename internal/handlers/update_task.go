package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/vrazinsky/go-final-project/internal/models"
	"github.com/vrazinsky/go-final-project/internal/nextdate"
	"github.com/vrazinsky/go-final-project/internal/utils"
)

func (h *Handlers) HandleUpdateTask(res http.ResponseWriter, req *http.Request) {
	var input models.TaskInput
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
	if err = json.Unmarshal(buf.Bytes(), &input); err != nil {
		logWriteErr(res.Write(ErrorResponse(err, "")))
		return
	}

	if len(input.Id) == 0 {
		logWriteErr(res.Write(ErrorResponse(nil, "no id")))
		return
	}
	_, err = strconv.ParseInt(input.Id, 10, 64)
	if err != nil {
		logWriteErr(res.Write(ErrorResponse(nil, "incorrect id")))
		return
	}
	err = input.Validate()
	if err != nil {
		logWriteErr(res.Write(ErrorResponse(err, "")))
		return
	}

	if input.Date == nil || *input.Date == "" {
		date = now
		dateToDb = now.Format(utils.Layout)
	} else {
		date, _ = time.Parse("20060102", *input.Date)
		dateToDb = *input.Date
	}
	if IsDateAfter(now, date) {
		if input.Repeat == nil || *input.Repeat == "" {
			dateToDb = now.Format(utils.Layout)
		} else {
			dateToDb, err = nextdate.NextDate(now, date.Format(utils.Layout), *input.Repeat)
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
		Id:      input.Id,
		Date:    dateToDb,
		Title:   input.Title,
		Comment: comment,
		Repeat:  repeat,
	}

	err = h.db.UpdateTask(task)
	if err != nil {
		logWriteErr(res.Write(ErrorResponse(err, "")))
		return
	}
	var response struct{}
	responseBytes, _ := json.Marshal(response)
	logWriteErr(res.Write(responseBytes))
}
