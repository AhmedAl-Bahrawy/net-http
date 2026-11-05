# Net/HTTP — Go Learning Journey ✅ COMPLETED

> **Status: Project Finished** 🎉
> This repository is the complete learning journey through Go's standard `net/http` library — from the first outbound HTTP request all the way to a full CRUD REST API backed by PostgreSQL. Next stop: **Chi** 🚂

---

## 📁 Projects Overview

### 🏠 Main (`main.go`)
Entrypoint that ties everything together:
- Loads environment variables from `.env` using **godotenv**
- Connects to PostgreSQL via a **pgxpool** connection pool
- Runs the client & request experiments (fetched once at startup)
- Registers server + REST API routes on one `http.ServeMux`
- Serves everything on `http://localhost:8080`

---

### 🌐 `client26-OCT/` — Simple HTTP Client

**Concepts:**
- Creating an `http.Client` and using `client.Get()` for simple GET requests
- Checking `response.StatusCode` (e.g., `http.StatusOK`)
- Closing response bodies with `defer res.Body.Close()`
- Decoding JSON responses with `json.NewDecoder().Decode()`

**API Used:** [Bored API](https://bored-api.appbrewery.com/random) — fetches a random activity suggestion.

---

### 📨 `request26-OCT/` — Building Custom Requests

**Concepts:**
- `http.NewRequest(method, url, body)` to build requests manually
- Explicit methods via `http.MethodGet`, `http.MethodPost`, ...
- Sending built requests with `client.Do(req)`
- **Key takeaway:** `client.Get()` = convenience, `client.Do(req)` = full control over headers, body, timeouts

---

### 🖥️ `server27-OCT/` — Servers & Handlers

**Concepts:**
- `http.HandleFunc("/", fn)` vs `http.Handle("/path", handler)`
- Implementing the `http.Handler` interface via `ServeHTTP(w, r)`
- Struct handlers that carry state (`greetHandler{Greeting: "Hi"}`)

**Routes:**
| Route | Handler Type | Description |
|-------|--------------|-------------|
| `/` | `http.HandlerFunc` | "Hello, This is Main Page" |
| `/greet` | `http.Handler` (struct) | Personalized greeting |

---

### ⚡ `restAPI27-OCT/` — Full CRUD REST API (PostgreSQL + pgx)

A production-style tasks API. Handlers are thin (parse input → call DB → write JSON), and all state lives in the database.

**Routes & Methods:**
| Method | Route | Description | Success Status |
|--------|-------|-------------|----------------|
| GET | `/tasks` | List all tasks | 200 |
| POST | `/tasks` | Create task (JSON body) | 201 Created |
| GET | `/tasks/{taskId}` | Get one task by ID | 200 |
| PATCH | `/tasks/{taskId}?title=...&completed=...` | Update task by ID | 200 |
| DELETE | `/tasks/{taskId}` | Delete task by ID | 200 |

**Data Model (`Task`):**
```go
type Task struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}
```

**Concepts covered:**
- Creating an isolated `http.NewServeMux()` instead of `DefaultServeMux`
- **Method-based routing** with `switch r.Method` inside a single path handler
- `{param}` path parameters + `r.PathValue("taskId")`
- Query parameters via `r.URL.Query().Get("completed")` + `strconv.ParseBool`
- Reading request JSON bodies: `json.NewDecoder(r.Body).Decode()`
- Writing statuses explicitly: `w.WriteHeader(http.StatusCreated)`
- Error responses with correct codes: `http.Error(w, msg, http.StatusBadRequest / NotFound / InternalServerError)`
- Rejecting unsupported verbs with `http.StatusMethodNotAllowed`
- Passing `r.Context()` down to database calls for cancellation support

**Database layer:**
- Explicit SQL with parameterized queries (`$1, $2, $3`) — no dynamic string concat
- `QueryRow(...).Scan(...)` for single rows, `Query(...)` + `rows.Next()` loop + `rows.Scan` for lists
- `INSERT ... RETURNING` to get the auto-generated ID back
- `result.RowsAffected() > 0` to detect "task not found" for PATCH / DELETE
- Handling `pgx.ErrNoRows` with `errors.Is`

---

### 🗄️ `database5-NOV/` — pgxpool Connection

**Concepts:**
- `pgxpool.New()` — pooled PostgreSQL connections (reuses connections across requests)
- Health-check via `pool.Ping()`
- Closing the pool cleanly with `defer pool.Close()`
- Reading the DSN from the `GOOSE_DBSTRING` environment variable (same var goose uses for migrations)

---

### 📜 `sql/migrations/` — goose Migrations

`20261105014940_create_tasks.sql`:
```sql
-- +goose Up
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE tasks;
```

**Concepts:**
- Timestamped, ordered migration files
- `-- +goose Up` (apply) / `-- +goose Down` (rollback) sections
- Managing schema evolution from the command line

---

## 🔑 Key net/http Concepts Learned

| Concept | Where |
|---------|-------|
| `http.Client` + `client.Get()` | client26-OCT |
| `http.NewRequest` + `client.Do(req)` | request26-OCT |
| `response.Body.Close()` cleanup | client26-OCT, request26-OCT |
| `json.NewDecoder().Decode()` | client/OCT request/REST API |
| `http.HandleFunc` / `http.Handle` / `ServeHTTP` | server27-OCT |
| `http.NewServeMux()` (custom mux) | restAPI27-OCT, main.go |
| Path parameters `{taskId}` + `r.PathValue()` | restAPI27-OCT |
| Query params + `r.URL.Query()` | restAPI27-OCT |
| Method-based routing via `r.Method` | restAPI27-OCT |
| `w.Header().Set()` / `w.WriteHeader()` | restAPI27-OCT |
| `http.Error()` with status codes | restAPI27-OCT |
| `r.Context()` for cancellation | restAPI27-OCT |
| Reading JSON from request bodies | restAPI27-OCT |
| Closures capturing dependencies in handlers | restAPI27-OCT |

**Bonus concepts beyond net/http:**
| Concept | Where |
|---------|-------|
| `pgxpool` connection pooling | database5-NOV |
| Parameterized SQL queries (`$1...`) | restAPI27-OCT |
| `INSERT ... RETURNING` | restAPI27-OCT |
| `RowsAffected()` for existence checks | restAPI27-OCT |
| `errors.Is(err, pgx.ErrNoRows)` | restAPI27-OCT |
| `.env` config with godotenv | main.go |
| Versioned SQL migrations (goose) | sql/migrations |

---

## 📌 Example Usage

```bash
# List all tasks
curl http://localhost:8080/tasks

# Create a task
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Chi", "completed": false}'

# Get one task
curl http://localhost:8080/tasks/1

# Update a task
curl -X PATCH "http://localhost:8080/tasks/1?title=Example&completed=true"

# Delete a task
curl -X DELETE http://localhost:8080/tasks/1
```

---

## 🏃 How to Run

1. **Create a `.env` file** (git-ignored) with your Postgres DSN:
   ```
   GOOSE_DBSTRING=postgres://user:password@localhost:5432/nethttp
   ```
2. **Apply migrations:**
   ```bash
   goose up
   ```
3. **Run the app:**
   ```bash
   go mod tidy
   go run main.go
   ```
4. Visit:
   - http://localhost:8080/ and /greet — classic server pages
   - http://localhost:8080/tasks — the REST API

---

## 🎯 What's Next

The `net/http` journey ends here. The next phase is **[Chi Router](https://github.com/go-chi/chi)**, which builds on the same standard-library foundations learned here and adds:
- Middleware chaining
- Composable sub-routers
- URL params with wildcard patterns (`r.URLParam`)
- Cleaner REST routing (`r.Get`, `r.Post`, `r.Patch`, `r.Delete`)

Learning Chi will feel natural: everything familiar from this repo (handlers, `ResponseWriter`, `Request`) carries over — Chi just organizes it better.
