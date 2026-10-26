package client26oct

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Activity struct {
    Activity string `json:"activity"`
}

func NewClient(url string) string{

	client := &http.Client{}

	res, err := client.Get(url)

	if err != nil {
		fmt.Println("Request Error: ", err)
		return err.Error()
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("Unexpected Status: ", res.Status,
					"\nCode: ", res.StatusCode)
		return string(res.StatusCode)
	}

	var acData Activity

	err = json.NewDecoder(res.Body).Decode(&acData)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	activity := "Hey We Found A Good Activity For You: " + acData.Activity
	return activity
}