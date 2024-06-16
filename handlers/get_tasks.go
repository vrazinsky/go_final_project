package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/vrazinsky/go-final-project/models"
)

func (h *Handlers) HandleGetTasks(res http.ResponseWriter, req *http.Request) {
	searchValue := req.FormValue("search")
	var (
		filterByDate  bool = false
		filterByTitle bool = false
	)
	if searchValue != "" {
		hasDate, _ := regexp.MatchString("^\\d{2}\\.\\d{2}\\.\\d{4}$", searchValue)
		if hasDate {
			date, err := time.Parse("02.01.2006", searchValue)
			if err != nil {
				res.Write(ErrorGetTasksResponse(err, ""))
				return
			}
			filterByDate = true
			searchValue = date.Format("20060102")
		} else {
			filterByTitle = true
			searchValue = fmt.Sprintf("%%%v%%", searchValue)
		}
	}
	rows, err := h.db.QueryContext(h.ctx, getTasksQuery,
		sql.Named("filterByTitle", filterByTitle),
		sql.Named("filterByDate", filterByDate),
		sql.Named("searchValue", searchValue))
	if err != nil {
		res.Write(ErrorGetTasksResponse(err, ""))
		return
	}
	defer rows.Close()
	tasks := make([]models.Task, 0)
	for rows.Next() {
		var task models.Task
		var taskId int
		err := rows.Scan(&taskId, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		task.Id = strconv.Itoa(taskId)
		if err != nil {
			res.Write(ErrorGetTasksResponse(err, ""))
			return
		}
		tasks = append(tasks, task)
	}
	response := models.GetTasksResponse{Tasks: &tasks}
	data, _ := json.Marshal(response)
	res.Write(data)
}
