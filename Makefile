#SHELL :=/bin/zsh
#export GOPATH:=/Users/anatoly.saukhin/GO
#export GOBIN:=${GOPATH}/bin
#export PATH:=${GOBIN}:${PATH}


.PHONY: build
build:
	go build -o bookings ./cmd/web/.

.PHONY: run
run: build
	 ./bookings

.PHONY: install
install:
	go install github.com/onsi/ginkgo/v2/ginkgo@latest

.PHONY: test
test:
	ginkgo ./...
