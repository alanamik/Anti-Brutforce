BIN := "./antibrutforce"
DOCKER_IMG="antibrutforce:latest"
DOCKER_IMG_TEST="antibrutforce:test"

build:
	go build -v -o $(BIN) ./cmd/cmd.go

run: build
	$(BIN) -config ./configs/dev.yml 

run-test: build
	$(BIN) -config ./configs/test.yml &

test:
	go test -timeout=90s -count=1 -v ./internal/...

test-http: docker-run-test
	go test -timeout=3m -count=1 -v ./tests/... 

docker-build:
	docker build -t $(DOCKER_IMG) .

docker-run:
	docker run $(DOCKER_IMG)

docker-build-test:
	docker build -t $(DOCKER_IMG_TEST) .

docker-run-test:
	docker run -d -p "8000:8000" $(DOCKER_IMG_TEST)  


install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.58.1

lint: install-lint-deps
	golangci-lint run ./...
