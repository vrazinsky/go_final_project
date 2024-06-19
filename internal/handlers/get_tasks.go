package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/vrazinsky/go-final-project/internal/models"
	"github.com/vrazinsky/go-final-project/internal/utils"
)

func (h *Handlers) HandleGetTasks(res http.ResponseWriter, req *http.Request) {
	searchValue := req.FormValue("search")
	var (
		filterByDate bool = false
		filterByText bool = false
	)
	if searchValue != "" {
		hasDate, _ := regexp.MatchString("^\\d{2}\\.\\d{2}\\.\\d{4}$", searchValue)
		if hasDate {
			date, err := time.Parse("02.01.2006", searchValue)
			if err != nil {
				logWriteErr(res.Write(ErrorGetTasksResponse(err, "")))
				return
			}
			filterByDate = true
			searchValue = date.Format(utils.Layout)
		} else {
			filterByText = true
			searchValue = fmt.Sprintf("%%%v%%", searchValue)
		}
	}
	tasks, err := h.db.GetTasks(filterByText, filterByDate, searchValue)
	if err != nil {
		logWriteErr(res.Write(ErrorGetTasksResponse(err, "")))
		return
	}
	response := models.GetTasksResponse{Tasks: &tasks}
	data, _ := json.Marshal(response)
	logWriteErr(res.Write(data))
}
