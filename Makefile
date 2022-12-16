# INTERNAL VARIABLES
#
# TARGETS FOR BUILD
#


build-node:
	@echo "[building node]"
	@docker build -t $(CHAIN)-node -f ./$(CHAIN)/docker/Dockerfile --progress tty .
	@echo "Build $(CHAIN) node docker image done ✔︎"

build-node-pristine:
	@echo "[building node]"
	@docker build --no-cache -t $(CHAIN)-node -f ./$(CHAIN)/docker/Dockerfile --progress tty .
	@echo "Build $(CHAIN) node docker image done ✔︎"

compose-node-up:
	@echo "[composing node up]"
	@docker-compose -f $(CHAIN)/docker-compose.yml up

compose-node-down:
	@echo "[composing node down]"
	@docker-compose -f $(CHAIN)/docker-compose.yml down

build-node-runner:
	@echo "[build node runner]"
	@go build -o bin/$(CHAIN)/runner src/cmd/$(CHAIN)/main.go
	@echo "Build $(CHAIN) node done ✔︎"

run-node-local:
	@echo "[run node local]"
	@export $$(cat $(CHAIN)/node.env) && nodemon --exec go run src/cmd/$(CHAIN)/main.go
