## 1st stage: build golang ethereum runner
FROM golang:alpine as builder

WORKDIR /usr/src/app

COPY ../../. .

RUN go build -o ethereum-runner src/cmd/ethereum/main.go

## 2nd stage: prepare container to run node
FROM node

WORKDIR /home/node

ARG NODE_URL
ENV NODE_URL ${NODE_URL}

COPY --from=builder /usr/src/app/ethereum-runner /

RUN npm install -g ganache

ENTRYPOINT ["/ethereum-runner"]
