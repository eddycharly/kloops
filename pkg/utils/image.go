package utils

import (
	"fmt"
	"net/http"
	"strconv"
)

// ImageTooBig checks if image is bigger than github limits
func ImageTooBig(url string) (bool, error) {
	// limit is 10MB
	limit := 10000000
	// try to get the image size from Content-Length header
	resp, err := http.Head(url) // #nosec
	if err != nil {
		return true, fmt.Errorf("HEAD error: %v", err)
	}
	if sc := resp.StatusCode; sc != http.StatusOK {
		return true, fmt.Errorf("failing %d response", sc)
	}
	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	if size > limit {
		return true, nil
	}
	return false, nil
}
