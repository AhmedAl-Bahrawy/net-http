# Net/HTTP — Go Learning Journey

This repository contains hands-on projects exploring Go's standard `net/http` library, progressing toward more advanced frameworks like **Chi**.

---

## 📁 Projects Overview

### 🏠 Main (`main.go`)
The entry point that demonstrates how to run both HTTP servers together on port 8080. It calls into each sub-package showing them in action simultaneously.

---

### 🌐 `client26-OCT/` — Making HTTP Requests with `http.Client`

**Concepts:**
- Creating an `http.Client` instance for making outbound HTTP requests
- Using `client.Get()` for simple GET requests
- Checking `response.StatusCode` (e.g., `http.StatusOK`)
- Properly closing response bodies with `defer res.Body.Close()`
- Decoding JSON responses using `json.NewDecoder()` and `Decode()`
- Basic error handling for request and decode failures

**API Used:** [Bored API](https://bored-api.appbrewery.com/random) — fetches random activity suggestions.

---

### 📨 `request26-OCT/` — Constructing HTTP Requests with `http.NewRequest`

**Concepts:**
- Using `http.NewRequest(method, url, body)` to build custom requests
- Specifying the HTTP method explicitly via `http.MethodGet`, `http.MethodPost`, etc.
- Passing headers and payloads (body) through the Request struct
- Sending built requests using `client.Do(req)` instead of the shortcut `.Get()`
- Difference between `client.Get()` (simple) vs `client.Do(req)` (full control)

**Key Takeaway:** This gives you full control over headers, methods, timeouts, and body content — essential for POST, PUT, DELETE, and authenticated requests.

---

### 🖥️ `server27-OCT/` — Building HTTP Servers & Handlers

**Concepts:**
- Registering handlers with `http.HandleFunc("/", handlerFunc)` — uses a pattern matcher
- Implementing the `http.Handler` interface by defining `ServeHTTP(w http.ResponseWriter, r *http.Request)`
- Difference between `http.HandlerFunc` (func-based) and `http.Handler` (struct-based, interface-based)
- Using custom structs as handlers to carry state (e.g., `greetHandler` holds a greeting string)
- Setting up routes with `http.Handle("/path", handler)`

**Routes:**
| Route | Handler | Description |
|-------|---------|-------------|
| `/` | `mainHandler` (func) | Prints "Hello, This is Main Page" |
| `/greet` | `greetHandler` (struct) | Prints a personalized greeting |

---

### ⚡ `restAPI27-OCT/` — REST API Server with Path Parameters

**Concepts:**
- Building a RESTful API endpoint pattern
- Using `{param}` path parameters in routes: `/tasks/{taskId}`
- Extracting path values with `r.PathValue("paramName")`
- Returning JSON responses with `json.NewEncoder(w).Encode()`
- Setting response headers (`Content-Type: application/json`)
- Returning error responses via `http.Error(w, message, statusCode)`
- Using closures with `http.HandlerFunc` to capture data (tasks slice) in handlers

**Routes:**
| Method | Route | Description |
|--------|-------|-------------|
| GET | `/tasks` | Returns all tasks as JSON |
| GET | `/tasks/{taskId}` | Returns a specific task by ID |

**Data Model (`Task`):**
```go
type Task struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}
```

---

## 🔑 Key net/http Concepts Learned

| Concept | Where |
|---------|-------|
| `http.Client{}` — creating an HTTP client | client26-OCT, request26-OCT |
| `client.Get(url)` — simple GET request | client26-OCT |
| `http.NewRequest()` — building custom requests | request26-OCT |
| `client.Do(req)` — sending any HTTP request | request26-OCT |
| `response.Body.Close()` — resource cleanup | client26-OCT, request26-OCT |
| `json.NewDecoder(res.Body).Decode()` — reading JSON | client26-OCT, request26-OCT |
| `http.HandleFunc()` — registering func handlers | server27-OCT |
| `http.Handle()` — registering interface handlers | server27-OCT |
| `ServeHTTP(w, r)` — implementing `http.Handler` | server27-OCT |
| `w.Header().Set()` — setting response headers | restAPI27-OCT |
| `json.NewEncoder(w).Encode()` — writing JSON responses | restAPI27-OCT |
| `http.Error()` — returning error responses | restAPI27-OCT |
| `r.PathValue("key")` — extracting URL path params | restAPI27-OCT |
| Closures capturing data in handlers | restAPI27-OCT |

---

## 🏃 How to Run

```bash
go mod tidy
go run main.go
```

Then visit:
- **http://localhost:8080/** → Main page
- **http://localhost:8080/greet** → Custom greeting
- **http://localhost:8080/tasks** → All tasks (JSON)
- **http://localhost:8080/tasks/1** → Single task by ID (JSON)

---

## 🎯 Next Steps

Once comfortable with `net/http`, explore **[Chi Router](https://github.com/go-chi/chi)** — a lightweight, idiomatic router for Go microservices, which builds on top of the standard library but adds:
- Pattern matching with wildcards
- Middlewares
- Sub-routers
- Cleaner routing syntax
