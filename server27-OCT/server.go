package server27oct

import (
	"fmt"
	"net/http"
)

type GreetHandler struct {
	Greeting string
}

func (h GreetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s, visitor!", h.Greeting)
}

func HandleServer(g string) {
	handler := GreetHandler{Greeting: g}
	http.Handle("/greet", handler)
}
