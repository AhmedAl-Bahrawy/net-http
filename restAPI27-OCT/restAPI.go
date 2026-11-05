package restapi27oct

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Task Struct Start

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// Task Struct End

// Database Operations Start

func CreateTask(ctx context.Context, pool *pgxpool.Pool, title string, completed bool) (Task, error) {
	var task Task

	query := `
		INSERT INTO tasks (title, completed)
		VALUES ($1, $2)
		RETURNING id, title, completed
	`

	err := pool.QueryRow(ctx, query, title, completed).Scan(
		&task.ID,
		&task.Title,
		&task.Completed,
	)

	return task, err
}

func PatchTask(ctx context.Context, pool *pgxpool.Pool, title string, completed bool, id int) (bool, error) {
	query := `
		UPDATE tasks
		SET title = $1, completed = $2
		WHERE id = $3
	`

	result, err := pool.Exec(ctx, query, title, completed, id)
	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}

func DeleteTask(ctx context.Context, pool *pgxpool.Pool, id int) (bool, error) {
	query := `
		DELETE FROM tasks
		WHERE id = $1
	`

	result, err := pool.Exec(ctx, query, id)
	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}

func GetTasks(ctx context.Context, pool *pgxpool.Pool) ([]Task, error) {
	query := `
		SELECT id, title, completed
		FROM tasks
		ORDER BY id
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Completed,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []Task{}
	}

	return tasks, nil
}

func GetTask(ctx context.Context, pool *pgxpool.Pool, id int) (Task, bool, error) {
	var task Task

	query := `
		SELECT id, title, completed
		FROM tasks
		WHERE id = $1
	`

	err := pool.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Completed,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return Task{}, false, nil
	}

	if err != nil {
		return Task{}, false, err
	}

	return task, true, nil
}

// Database Operations End

// JSON / HTTP Handlers Start

func createTaskJson(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		createdTask, err := CreateTask(
			r.Context(),
			pool,
			newTask.Title,
			newTask.Completed,
		)

		if err != nil {
			http.Error(w, "Failed to create task", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(createdTask)
	}
}

func patchTaskJson(pool *pgxpool.Pool) http.HandlerFunc {
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

		found, err := PatchTask(
			r.Context(),
			pool,
			title,
			completed,
			taskId,
		)

		if err != nil {
			http.Error(w, "Failed to update task", http.StatusInternalServerError)
			return
		}

		if !found {
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

func deleteTaskJson(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("taskId")
		taskId, err := strconv.Atoi(id)

		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		found, err := DeleteTask(r.Context(), pool, taskId)

		if err != nil {
			http.Error(w, "Failed to delete task", http.StatusInternalServerError)
			return
		}

		if !found {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Task deleted successfully",
		})
	}
}

func getTasksJson(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := GetTasks(r.Context(), pool)

		if err != nil {
			http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(tasks)
	}
}

func getTaskJson(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("taskId")
		taskId, err := strconv.Atoi(id)

		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		task, found, err := GetTask(r.Context(), pool, taskId)

		if err != nil {
			http.Error(w, "Failed to retrieve task", http.StatusInternalServerError)
			return
		}

		if !found {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(task)
	}
}

// JSON / HTTP Handlers End

// Endpoints Start

func tasksHandler(mux *http.ServeMux, pool *pgxpool.Pool) {
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTasksJson(pool)(w, r)

		case http.MethodPost:
			createTaskJson(pool)(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func taskHandler(mux *http.ServeMux, pool *pgxpool.Pool) {
	mux.HandleFunc("/tasks/{taskId}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTaskJson(pool)(w, r)

		case http.MethodPatch:
			patchTaskJson(pool)(w, r)

		case http.MethodDelete:
			deleteTaskJson(pool)(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

// Endpoints End

func CreateRestAPIServer(pool *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()

	tasksHandler(mux, pool)
	taskHandler(mux, pool)

	return mux
}
