
# KLoops Dashboard backend

## Routes

The routes below are the routes supported by the dashboard backend.

They are evaluated in order, the first match serves the request.

### /proxy/{rest:.*}

Proxy to the kubernetes api server

This route supports all methods.

### /api/pluginhelp

Get plugins help

This route supports the following methods:

- GET


### /api/plugins

List plugin configs

This route supports the following methods:

- GET


### /api/plugins/{name}

Get plugin config

This route supports the following methods:

- GET


### /api/repos

List repo configs

This route supports the following methods:

- GET


### /api/repos/{name}

Get repo config

This route supports the following methods:

- GET


### /api/repos/{name}

Create repo config

This route supports the following methods:

- POST


### /api/hooks/{name}

Setup repo config hooks

This route supports the following methods:

- POST


### /socket

Connect to websocket

This route supports the following methods:

- GET


### / (prefix)

Static content

This route supports the following methods:

- GET

