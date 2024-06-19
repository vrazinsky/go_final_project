package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func (h *Handlers) HandleGetTask(res http.ResponseWriter, req *http.Request) {
	idStr := req.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logWriteErr(res.Write(ErrorResponse(nil, "incorrect input data")))
		return
	}
	task, err := h.db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logWriteErr(res.Write(ErrorResponse(nil, "task not found")))
		} else {
			logWriteErr(res.Write(ErrorResponse(err, "")))
		}
		return
	}
	data, _ := json.Marshal(task)
	logWriteErr(res.Write(data))
}
