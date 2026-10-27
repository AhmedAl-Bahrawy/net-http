package main

import (
	"fmt"
	client26oct "net-http/client26-OCT"
	request26oct "net-http/request26-OCT"
	server27oct "net-http/server27-OCT"
)

func main() {
	url26OCT := "https://bored-api.appbrewery.com/random"

	client := client26oct.NewClient(url26OCT)
	fmt.Println(client)

	request := request26oct.NewRequest(url26OCT)
	fmt.Println(request)

	server27oct.HandleServer("Hi")
}
