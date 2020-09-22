module github.com/eddycharly/kloops

go 1.13

require (
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/go-logr/logr v0.1.0
	github.com/gorilla/mux v1.8.0
	github.com/jenkins-x/go-scm v1.5.166
	github.com/nbio/st v0.0.0-20140626010706-e9e8d9816f32 // indirect
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.9.1
	gopkg.in/yaml.v1 v1.0.0-20140924161607-9f9df34309c0
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	sigs.k8s.io/controller-runtime v0.5.0
	sigs.k8s.io/yaml v1.1.0
)

// replace github.com/jenkins-x/go-scm => /Users/charlesbreteche/Documents/dev/eddycharly/go-scm
