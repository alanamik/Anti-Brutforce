BIN := "./antibrutforce"
DOCKER_IMG="antibrutforce:latest"

build:
	go build -v -o $(BIN) ./cmd/cmd.go

run: build
	$(BIN) -config ./configs/dev.yml 

test:
	go test -timeout=10s -count=1 -v ./internal/...
#./pkg/...


docker-build:
	docker build -t $(DOCKER_IMG) .

docker-run:
	docker run $(DOCKER_IMG)


install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.58.1

lint: install-lint-deps
	golangci-lint run ./...
