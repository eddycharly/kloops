package utils

import "github.com/jenkins-x/go-scm/scm"

func QuoteAuthorForComment(client *scm.Client, author string) string {
	if client.Driver == scm.DriverStash {
		return `"` + author + `"`
	}
	return author
}
