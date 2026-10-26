package main

import (
	"fmt"
	client26oct "net-http/client26-OCT"
)

func main(){
	url26OCT := "https://bored-api.appbrewery.com/random"
	client := client26oct.NewClient(url26OCT)
	fmt.Println(client)
}