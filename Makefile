ENABLE_VENDOR=1
export GOPATH=$(shell echo $$GOPATH:`pwd`)
export GO15VENDOREXPERIMENT=${ENABLE_VENDOR}
GO_PKGS=$(shell go list ./... | grep -v '/vendor/')

default: build

build:
	go build

dep:
	godep get -t ${GO_PKGS}
	godep save -t ${GO_PKGS}

.PHONY: witch
