package handlers

import (
	"net/http"
	"time"

	"github.com/vrazinsky/go-final-project/nextdate"
)

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
	nextDate, err := nextdate.NextDate(nowDate, date, repeat)
	if err != nil {
		logWriteErr(res.Write([]byte(err.Error())))
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	logWriteErr(res.Write([]byte(nextDate)))
}
