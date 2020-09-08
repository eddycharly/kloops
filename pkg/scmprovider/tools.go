package scmprovider

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jenkins-x/go-scm/scm"
)

type Tools struct {
	client *scm.Client
}

func (s Tools) QuoteAuthorForComment(author string) string {
	if s.client.Driver == scm.DriverStash {
		return `"` + author + `"`
	}
	return author
}

// ImageTooBig checks if image is bigger than github limits
func (s Tools) ImageTooBig(url string) (bool, error) {
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
