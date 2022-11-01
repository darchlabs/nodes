## 1st stage: build golang ethereum runner
FROM golang:alpine as builder

WORKDIR /usr/src/app

COPY ../../. .

RUN go build -o eth-runner src/cmd/ethereum/main.go

## 2nd stage: prepare container to run node
FROM node

WORKDIR /home/node

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
## NODE_URL
ARG NODE_URL
ENV NODE_URL ${NODE_URL}
## BLOCK_NUMBER
ARG BLOCK_NUMBER
ENV BLOCK_NUMBER ${BLOCK_NUMBER}



COPY --from=builder /usr/src/app/eth-runner /home/node

RUN npm install -g ganache
RUN echo "setInterval(function () {console.log('Every 5 secondes');}, 5000);console.log('now');" >> index.js

#CMD ["./eth-runner"]
CMD ["node", "index.js"]
