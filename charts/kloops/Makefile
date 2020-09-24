GO_MOD := GO111MODULE=on go

ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

HELMDOCS := $(GOBIN)/helm-docs
$(HELMDOCS):
	pushd /tmp; $(GO_MOD) get -u github.com/norwoodj/helm-docs/cmd/helm-docs@v0.15.0; popd

lint:
	helm lint

docs: $(HELMDOCS)
	$(HELMDOCS)

build: lint docs
