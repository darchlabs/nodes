## 1st stage: build golang ethereum runner
FROM golang as builder

WORKDIR /usr/src/app

COPY ../../. .

RUN go build -o nodes cmd/nodes/main.go

## 2nd stage: prepare container to run node
FROM node

WORKDIR /home/nodes

# Environment var
## ENVIRONMENT
ARG ENVIRONMENT
ENV ENVIRONMENT ${ENVIRONMENT}
## CHAIN
ARG CHAIN
ENV CHAIN ${CHAIN}
## API_SERVER_PORT
ARG API_SERVER_PORT
ENV API_SERVER_PORT ${API_SERVER_PORT}
## BASE_CHAIN_DATA_PATH
ARG BASE_CHAIN_DATA_PATH
ENV BASE_CHAIN_DATA_PATH ${BASE_CHAIN_DATA_PATH}

COPY --from=builder /usr/src/app/nodes /home/nodes

RUN npm install -g ganache

CMD ["./nodes"]
