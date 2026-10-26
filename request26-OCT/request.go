package request26oct

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Activity struct {
    Activity string `json:"activity"`
}

func NewRequest(url string) string{
	// Connect the Client
	client := &http.Client{}

	res, err := client.Get(url)

	if err != nil {
		fmt.Println("Request Error: ", err)
		return err.Error()
	}

	// Prepare for the Request
	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	// Error Handle
	if err != nil {
		fmt.Println("Handle Error: ", err)
		return err.Error()
	}

	// Post the Request and Wait for the Answer/Response
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Request Error: ", err)
		return err.Error()
	}
	defer res.Body.Close()

	// define Resonce Body
	var acData Activity

	// Decode the Response
	err = json.NewDecoder(res.Body).Decode(&acData)
	if err != nil {
		fmt.Println("Decode Error: ", err)
		return err.Error()
	}

	return acData.Activity
}