package unsplash

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const baseURL = "https://api.unsplash.com/photos/random"

type result struct {
	ID     string   `json:"id"`
	Images imageSet `json:"urls"`
}

type imageSet struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
	Small   string `json:"small"`
	Thumb   string `json:"thumb"`
}

// Search performs a search request, i can take a category and allows searching through gif images
func Search(query string, key string) (string, error) {
	resp, err := http.Get(getURL(query, key)) // #nosec
	if err != nil {
		return "", fmt.Errorf("could not read result: %v", err)
	}
	defer resp.Body.Close()
	if sc := resp.StatusCode; sc > 299 || sc < 200 {
		return "", fmt.Errorf("failing %d response", sc)
	}
	var res result
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	if res.Images.Small == "" {
		return "", fmt.Errorf("no image url in response")
	}
	return res.Images.Small, nil
}

func getURL(query string, key string) string {
	uri := baseURL
	uri += "?query=" + url.QueryEscape(query)
	if key != "" {
		uri += "&client_id=" + url.QueryEscape(key)
	}
	return uri
}
