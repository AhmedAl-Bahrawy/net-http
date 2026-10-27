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

func (t Task) createTask(title string, completed bool) {
	// I will write this function tommorow
}

func getTasksData() (tasks *[]Task) {
	// I will edit this function
	tasks = &[]Task{
		{ID: 1, Title: "Learn Go", Completed: false},
		{ID: 2, Title: "Build an API", Completed: false},
	}
	return
}

func getTasks(tasks *[]Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(*tasks)
	}
}

func getTask(tasks *[]Task) http.HandlerFunc {
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

	// need to edit how to import this tasks
	tasks := *getTasksData()

	http.HandleFunc("/tasks", getTasks(&tasks))
	http.HandleFunc("/tasks/{taskId}", getTask(&tasks))

}
