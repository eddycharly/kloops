# kloops

## Creating the project from scratch

1. Run `./install-kubebuilder.sh` to install [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder)
1. Run `go mod init github.com/eddycharly/kloops`
1. Run `kubebuilder init --domain kloops.io` to create the project

## Creating the apis

### RepoConfig

`kubebuilder create api --group config --version v1alpha1 --kind RepoConfig`

