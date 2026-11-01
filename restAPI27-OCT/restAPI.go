package restapi27oct

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type Tasks []Task

func (t *Tasks) GetTasks() []Task {
	return *t
}

func (t *Tasks) CreateTask(title string, completed bool) {
	task := Task{
		ID:        len(*t) + 1,
		Title:     title,
		Completed: completed,
	}
	*t = append(*t, task)
}

func getTasksJson(tasks *[]Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(*tasks)
	}
}

func getTaskJson(tasks *[]Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := r.PathValue("taskId")
		taskId, err := strconv.Atoi(id)

		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		for _, task := range *tasks {
			if task.ID == taskId {
				json.NewEncoder(w).Encode(task)
				return
			}
		}
		http.Error(w, "Task not found", http.StatusNotFound)
	}
}

func CreateRestAPIServer() {
	var tasks Tasks

	tasks.CreateTask("Learn Go", false)
	tasks.CreateTask("Build an API", false)
	tasks.CreateTask("Learn PostgreSQL", true)

	allTasks := tasks.GetTasks()

	http.HandleFunc("/tasks", getTasksJson(&allTasks))
	http.HandleFunc("/tasks/{taskId}", getTaskJson(&allTasks))

}
