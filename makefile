BIN := "./antibrutforce"
DOCKER_IMG="antibrutforce:latest"

build:
	go build -v -o $(BIN) ./cmd/cmd.go

run: build
	$(BIN) -config ./configs/dev.yml 

run-test: build
	$(BIN) -config ./configs/test.yml &

test:
	go test -timeout=90s -count=1 -v ./internal/...

test-http: 
	go test -timeout=3m -count=1 -v ./tests/...
#./pkg/...


docker-build:
	docker build -t $(DOCKER_IMG) .

docker-run:
	docker run $(DOCKER_IMG)


install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.58.1

lint: install-lint-deps
	golangci-lint run ./...
