# load .env file
include node.env
export $(shell sed 's/=.*//' node.env)

SERVICE_NAME=testing-nodes-docker
DOCKER_USER=darchlabs

build:
	@echo "[building node]"
	@docker build -t darchlabs/nodes -f ./Dockerfile --progress tty .
	@echo "Build darchlabs-nodes docker image done ✔︎"

build-pristine:
	@echo "[building node]"
	@docker build --no-cache -r darchlabs/nodes -f ./Dockerfile --progress tty .
	@echo "Build darchlabs-nodes docker image done ✔︎"

compose-up:
	@echo "[composing node up]"
	@docker-compose -f docker-compose.yml up

compose-down:
	@echo "[composing node down]"
	@docker-compose -f docker-compose.yml down

build-local:
	@echo "[build darchlabs-nodes local]"
	@go build -o bin/nodes/nodes cmd/nodes/main.go
	@echo "Build darchlabs-nodes done ✔︎"

dev:
	@echo "[run node local]"
	@export $$(cat node.env) && nodemon --exec go run cmd/nodes/main.go

docker-login:
	@echo "[docker] Login to docker..."
	@docker login -u $(DOCKER_USER) -p $(DOCKER_PASS)

docker: docker-login
	@echo "[docker] pushing $(REGISTRY_URL)/$(SERVICE_NAME):$(VERSION)"
	@docker buildx create --use
	@docker buildx build --platform linux/amd64,linux/arm64  --push -t $(DOCKER_USER)/$(SERVICE_NAME):$(VERSION)	.
