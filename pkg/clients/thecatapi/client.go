package thecatapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const baseURL = "https://api.unsplash.com/photos/random"

type catResult struct {
	Image string `json:"url"`
}

// Search performs a search request, i can take a category and allows searching through gif images
func Search(category string, movieCat bool, key string) (string, error) {
	resp, err := http.Get(getURL(category, movieCat, key)) // #nosec
	if err != nil {
		return "", fmt.Errorf("could not read cat: %v", err)
	}
	defer resp.Body.Close()
	if sc := resp.StatusCode; sc > 299 || sc < 200 {
		return "", fmt.Errorf("failing %d response", sc)
	}
	var cats []catResult
	if err = json.NewDecoder(resp.Body).Decode(&cats); err != nil {
		return "", err
	}
	if len(cats) < 1 {
		return "", errors.New("no cats in response")
	}
	cat := cats[0]
	if cat.Image == "" {
		return "", errors.New("no image url in response")
	}
	return cat.Image, nil
}

func getURL(category string, movieCat bool, key string) string {
	uri := baseURL
	if category != "" {
		uri += "&category=" + url.QueryEscape(category)
	}
	if key != "" {
		uri += "&api_key=" + url.QueryEscape(key)
	}
	if movieCat {
		uri += "&mime_types=gif"
	}
	return uri
}
