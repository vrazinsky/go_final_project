package models

type Task struct {
	Id      string  `json:"id"`
	Date    string  `json:"date"`
	Title   string  `json:"title"`
	Comment *string `json:"comment,omitempty"`
	Repeat  *string `json:"repeat,omitempty"`
}

type AddTaskInput struct {
	Date    *string `json:"date,omitempty"`
	Title   string  `json:"title"`
	Comment *string `json:"comment,omitempty"`
	Repeat  *string `json:"repeat,omitempty"`
}

type UpdateTaskInput struct {
	Id      string  `json:"id"`
	Date    *string `json:"date,omitempty"`
	Title   string  `json:"title"`
	Comment *string `json:"comment,omitempty"`
	Repeat  *string `json:"repeat,omitempty"`
}

type AddTaskResponse struct {
	Id    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type GetTasksResponse struct {
	Error string  `json:"error,omitempty"`
	Tasks *[]Task `json:"tasks"`
}

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

type SignInInput struct {
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}
