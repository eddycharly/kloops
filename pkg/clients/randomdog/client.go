package randomdog

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "https://random.dog/woof.json"

type dogResult struct {
	URL string `json:"url"`
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
		return "", fmt.Errorf("could not read dog from %s: %v", uri, err)
	}
	defer resp.Body.Close()
	var a dogResult
	if err = json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return "", err
	}
	return a.URL, nil
}
