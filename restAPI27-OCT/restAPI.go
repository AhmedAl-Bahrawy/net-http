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

func createTaskJson(tasks *Tasks) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid Http Method", http.StatusMethodNotAllowed)
			return
		}

		var newTask Task

		err := json.NewDecoder(r.Body).Decode(&newTask)

		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if newTask.Title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}

		tasks.CreateTask(newTask.Title, newTask.Completed)

		createdTask := (*tasks)[len(*tasks)-1]

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(createdTask)
	}
}

func getTasksJson(tasks *Tasks) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid http Method", http.StatusMethodNotAllowed)
			return
		}

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

func tasksHandler(tasks Tasks) {
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTasksJson(&tasks)(w, r)

		case http.MethodPost:
			createTaskJson(&tasks)(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

}

func CreateRestAPIServer() {
	var tasks Tasks

	tasks.CreateTask("Learn Go", false)
	tasks.CreateTask("Build an API", false)
	tasks.CreateTask("Learn PostgreSQL", true)

	taskSlices := tasks.GetTasks()

	tasksHandler(tasks)

	http.HandleFunc("/tasks/{taskId}", getTaskJson(&taskSlices))

}
