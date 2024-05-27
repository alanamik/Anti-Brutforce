BIN := "./antibrutforce"
DOCKER_IMG="antibrutforce:latest"
GEN_DIR=internal/gen

swagger-gen:
	if ! [ -d $(GEN_DIR) ]; then \
	    mkdir $(GEN_DIR); \
	elif [ -d $(GEN_DIR) ]; then \
		rm -rf $(GEN_DIR); \
		mkdir $(GEN_DIR); \
	fi && \
	swagger generate server -t internal/gen -f ./api/swagger.yml --exclude-main -A anti-brutForce && \
	go mod tidy && \
	git add $(GEN_DIR)

run: build
	$(BIN) -config ./configs/dev.yml 

build:
	go build -v -o $(BIN) ./cmd/cmd.go

docker-compose-run:
	docker compose up -d

docker-build:
	docker build -t $(DOCKER_IMG) .

docker-run:
	docker run $(DOCKER_IMG)
