package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main(){
	client := &http.Client{}

	res, err := client.Get("https://bored-api.appbrewery.com/random")

	if err != nil {
		fmt.Println("Request Error: ", err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("Unexpected Status: ", res.Status,
					"\nCode: ", res.StatusCode)
		return
	}

	var data map[string]any

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Hey We Found A Good Activity For U: \n%v\n", data["activity"])
}