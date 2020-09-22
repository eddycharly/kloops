package scmprovider

import (
	"strings"

	"github.com/jenkins-x/go-scm/scm"
)

type Client struct {
	client       *scm.Client
	Issues       Issues
	PullRequests PullRequests
	Repositories Repositories
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
		Repositories: Repositories{
			client: client.Repositories,
		},
		Tools: Tools{
			client: client,
		},
	}
}

// NormLogin normalizes GitHub login strings
var NormLogin = strings.ToLower
