package main

import (
	"fmt"
	"log"
	"net/http"

	// Go maps this prefix locally because it matches your go.mod module declaration
	client26oct "github.com/AhmedAl-Bahrawy/net-http/client26-OCT"
	request26oct "github.com/AhmedAl-Bahrawy/net-http/request26-OCT"
	restapi27oct "github.com/AhmedAl-Bahrawy/net-http/restAPI27-OCT"
	server27oct "github.com/AhmedAl-Bahrawy/net-http/server27-OCT"
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
