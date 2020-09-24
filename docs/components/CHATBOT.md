# KLoops components

The chatbot component is responsible of interacting with an scm provider (GitHub, gitlab, Gitea, etc...).

This works by receiving and reeacting to webhooks sent by those scm providers when events occur in a repositorry (git pushes, pull requests, comments, reviews, etc...).

The chatbot is composed of several plugins, each plugin serving a different purpose.

The available plugins are documented here (TODO).

## Endpoints

The chatbot component exposes one endpoint per scm provider. Assuming it is deployed to the `example.com` domain, the following endpoint will be used:

| SCM provider      | Endpoint                      |
|-------------------|-------------------------------|
| GitHub            | `example.com/hook/github`     |
| Gitlab            | `example.com/hook/gitlab`     |
| Gitea             | `example.com/hook/gitea`      |
| Bitbucket cloud   | `example.com/hook/bitbucket`  |
| Bitbucket server  | `example.com/hook/stash`      |
| Gogs              | `example.com/hook/gogs`       |
