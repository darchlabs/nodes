# 1st stage: clone gnoland --------------
FROM ubuntu as gno

WORKDIR home/gnolang

RUN apt-get -y update

RUN apt-get -y install git

RUN git clone https://github.com/gnolang/gno.git

# 2nd stage: build gnoland --------------
FROM golang as gnobuilder

WORKDIR /usr/src/gno

COPY --from=gno home/gnolang/gno/. /usr/src/gno/.

#RUN go mod download

RUN go build -o build/gnoland cmd/gnoland/main.go

## 3rd stage: build golang ethereum runner --------------
FROM golang as builder

WORKDIR /usr/src/app

COPY ../../. .

RUN go build -o eth-runner src/cmd/ethereum/main.go

CMD ["./eth-runner"]

# 4th stage: prepare container to run node --------
FROM node

WORKDIR /home/node

## Environment var
### ENVIRONMENT
ARG ENVIRONMENT
ENV ENVIRONMENT ${ENVIRONMENT}
# CHAIN
ARG CHAIN
ENV CHAIN ${CHAIN}
# API_SERVER_PORT
ARG API_SERVER_PORT
ENV API_SERVER_PORT ${API_SERVER_PORT}
# BASE_CHAIN_DATA_PATH
ARG BASE_CHAIN_DATA_PATH
ENV BASE_CHAIN_DATA_PATH ${BASE_CHAIN_DATA_PATH}
# NODE_URL
ARG NODE_URL
ENV NODE_URL ${NODE_URL}
# BLOCK_NUMBER
ARG BLOCK_NUMBER
ENV BLOCK_NUMBER ${BLOCK_NUMBER}

COPY --from=builder /usr/src/app/eth-runner /home/node

# --------- GNOLAND ------------
COPY --from=gnobuilder /usr/src/gno/. /usr/src/gno/.

RUN npm install -g ganache

CMD ["./eth-runner"]

