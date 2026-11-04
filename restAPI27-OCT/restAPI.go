package restapi27oct

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Task / Tasks Struct Start

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type Tasks []Task

func (t *Tasks) CreateTask(title string, completed bool) {
	task := Task{
		ID:        len(*t) + 1,
		Title:     title,
		Completed: completed,
	}
	*t = append(*t, task)
}

func (t *Tasks) PatchTask(title string, completed bool, id int) bool {
	for i := range *t {
		if (*t)[i].ID == id {
			(*t)[i].Title = title
			(*t)[i].Completed = completed
			return true
		}
	}
	return false
}

// Task / Tasks Struct End

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

func patchTaskJson(tasks *Tasks) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()

		title := queryParams.Get("title")
		if title == "" {
			http.Error(w, "Invalid string value for 'title'", http.StatusBadRequest)
			return
		}

		completed, err := strconv.ParseBool(queryParams.Get("completed"))
		if err != nil {
			http.Error(w, "Invalid boolean value for 'completed'", http.StatusBadRequest)
			return
		}

		id := r.PathValue("taskId")
		taskId, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		found := tasks.PatchTask(title, completed, taskId)

		if found == false {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Task updated successfully",
		})
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

func getTaskJson(tasks *Tasks) http.HandlerFunc {
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

func tasksHandler(tasks *Tasks) {
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTasksJson(tasks)(w, r)

		case http.MethodPost:
			createTaskJson(tasks)(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/create", createTaskJson(tasks))
}

func taskHandler(tasks *Tasks) {
	http.HandleFunc("/tasks/{taskId}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTaskJson(tasks)(w, r)

		case http.MethodPatch:
			patchTaskJson(tasks)(w, r)

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

	tasksHandler(&tasks)
	taskHandler(&tasks)
}
