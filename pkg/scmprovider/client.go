package scmprovider

import "github.com/jenkins-x/go-scm/scm"

type Client struct {
	client       *scm.Client
	Issues       Issues
	PullRequests PullRequests
	Git          Git
	Tools        Tools
}

func NewClient(client *scm.Client) Client {
	return Client{
		client: client,
		Issues: Issues{
			client: client.Issues,
		},
		PullRequests: PullRequests{
			client: client.PullRequests,
		},
		Git: Git{
			client: client.Git,
		},
		Tools: Tools{
			client: client,
		},
	}
}
