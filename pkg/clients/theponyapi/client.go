package theponyapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const baseURL = "https://theponyapi.com/api/v1/pony/random"

type result struct {
	Pony pony `json:"pony"`
}

type pony struct {
	Representations representations `json:"representations"`
}

type representations struct {
	Full  string `json:"full"`
	Small string `json:"small"`
}

// Search performs a search request, it can take a queery
func Search(tags string) (string, error) {
	uri := string(getURL(tags))
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
	var a result
	if err = json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return "", err
	}
	return a.Pony.Representations.Small, nil
}

func getURL(tags string) string {
	uri := baseURL
	uri += "?q=" + url.QueryEscape(tags)
	return uri
}
