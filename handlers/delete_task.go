package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *Handlers) HandleDeleteTask(res http.ResponseWriter, req *http.Request) {
	idStr := req.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logWriteErr(res.Write(ErrorResponse(nil, "incorrect input data")))
		return
	}
	_, err = h.db.ExecContext(h.ctx, deleteTaskQuery, sql.Named("id", id))
	if err != nil {
		logWriteErr(res.Write(ErrorResponse(err, "")))
		return
	}
	var response struct{}
	responseBytes, _ := json.Marshal(response)
	logWriteErr(res.Write(responseBytes))
}
