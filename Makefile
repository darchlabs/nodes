# INTERNAL VARIABLES
#
# TARGETS FOR BUILD
#


build-node:
	@echo "[building node]"
	@docker build --no-cache -t $(CHAIN)-node -f ./$(CHAIN)/docker/Dockerfile --progress tty .

compose-node-up:
	@echo "[composing node up]"
	@docker-compose -f $(CHAIN)/docker-compose.yml up

compose-node-down:
	@echo "[composing node down]"
	@docker-compose -f $(CHAIN)/docker-compose.yml down

