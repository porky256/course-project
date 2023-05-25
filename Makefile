#export GOPATH:=/Users/anatoly.saukhin/GO
#export GOBIN:=${GOPATH}/bin
#export PATH:=${GOBIN}:${PATH}
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=2341
export POSTGRES_DB=db

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
	go test -failfast ./...

.PHONY: test-coverage
test-coverage:
	ginkgo -r -v -race --trace  --coverprofile=.coverage-report.out ./...

.PHONY: show-test-coverage
show-test-coverage: test-coverage
	go tool cover -html=.coverage-report.out

.PHONY: run-db
run-db:
	docker-compose up db

.PHONY: stop-db
stop-db:
	docker-compose stop

migration-up:
	printenv
	migrate -path data/migrations -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@0.0.0.0:5432/$(POSTGRES_DB)?sslmode=disable' up

migration-down:
	migrate -path data/migrations -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@0.0.0.0:5432/$(POSTGRES_DB)?sslmode=disable' down

create-new-migration:
	migrate create -ext .sql -dir data/migrations seed-restrictions