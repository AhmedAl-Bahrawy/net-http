package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	// Go maps this prefix locally because it matches your go.mod module declaration
	client26oct "github.com/AhmedAl-Bahrawy/net-http/client26-OCT"
	request26oct "github.com/AhmedAl-Bahrawy/net-http/request26-OCT"
	restapi27oct "github.com/AhmedAl-Bahrawy/net-http/restAPI27-OCT"
	server27oct "github.com/AhmedAl-Bahrawy/net-http/server27-OCT"

	database5nov "github.com/AhmedAl-Bahrawy/net-http/database5-NOV"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	pool, err := database5nov.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	url26OCT := "https://bored-api.appbrewery.com/random"

	client := client26oct.NewClient(url26OCT)
	fmt.Println(client)

	request := request26oct.NewRequest(url26OCT)
	fmt.Println(request)

	server27oct.HandleServer("Hi")

	mux := restapi27oct.CreateRestAPIServer(pool)

	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
