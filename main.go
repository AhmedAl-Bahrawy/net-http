package main

import (
	"fmt"
	"log"
	client26oct "net-http/client26-OCT"
	request26oct "net-http/request26-OCT"
	restapi27oct "net-http/restAPI27-OCT"
	server27oct "net-http/server27-OCT"
	"net/http"
)

func main() {
	url26OCT := "https://bored-api.appbrewery.com/random"

	client := client26oct.NewClient(url26OCT)
	fmt.Println(client)

	request := request26oct.NewRequest(url26OCT)
	fmt.Println(request)

	server27oct.HandleServer("Hi")

	restapi27oct.CreateRestAPIServer()

	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
