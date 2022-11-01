## 1st stage: build golang ethereum runner
FROM golang:alpine as builder

WORKDIR /usr/src/app

COPY ../../. .

RUN go build -o eth-runner src/cmd/ethereum/main.go

## 2nd stage: prepare container to run node
FROM node

WORKDIR /home/node

ARG NODE_URL
ENV NODE_URL ${NODE_URL}
ARG BLOCK_NUMBER
ENV BLOCK_NUMBER ${BLOCK_NUMBER}

COPY --from=builder /usr/src/app/eth-runner /home/node

RUN npm install -g ganache
RUN echo "setInterval(function () {console.log('Every 5 secondes');}, 5000);console.log('now');" >> index.js

#CMD ["./eth-runner"]
CMD ["node", "index.js"]
