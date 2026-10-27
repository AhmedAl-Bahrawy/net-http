package server27oct

import (
	"fmt"
	"net/http"
)

type greetHandler struct {
	Greeting string
}

func (h greetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s, visitor!", h.Greeting)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, This is Main Page")
}

func HandleServer(g string) {
	http.HandleFunc("/", mainHandler)

	gHandler := greetHandler{Greeting: g}
	http.Handle("/greet", gHandler)

	http.ListenAndServe(":8080", nil)
}
