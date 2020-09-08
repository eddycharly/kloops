package icanhazdadjoke

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "https://icanhazdadjoke.com"

type result struct {
	Joke string `json:"joke"`
}

// Get returns a ranom dog
func Get() (string, error) {
	uri := string(baseURL)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", fmt.Errorf("could not create request %s: %v", uri, err)
	}
	req.Header.Add("Accept", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not read joke from %s: %v", uri, err)
	}
	defer resp.Body.Close()
	var a result
	if err = json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return "", err
	}
	if a.Joke == "" {
		return "", fmt.Errorf("result from %s did not contain a joke", uri)
	}
	return a.Joke, nil
}
