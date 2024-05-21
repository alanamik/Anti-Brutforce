BIN := "./antibrutforce"
DOCKER_IMG="--"
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

run-redis:
	docker compose up -d
