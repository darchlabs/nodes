# INTERNAL VARIABLES
#
# TARGETS FOR BUILD
#


build-darch-node:
	@echo "[build-darch]"
	@docker build --no-cache -t darch-node -f ./darch-node/docker/Dockerfile --progress tty .

compose-darch:
	@echo "[compose-darch]"

