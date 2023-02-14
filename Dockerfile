## 1st stage: build golang ethereum runner
FROM golang as builder

WORKDIR /usr/src/app

COPY . .

RUN go build -o nodes cmd/nodes/main.go

## 2nd stage: prepare container to run node
FROM golang

WORKDIR /home/nodes

## ENVIRONMENT
ARG ENVIRONMENT
ENV ENVIRONMENT ${ENVIRONMENT}
## API_SERVER_PORT
ARG API_SERVER_PORT
ENV API_SERVER_PORT ${API_SERVER_PORT}
## REDIS_URL
ARG REDIS_URL
ENV REDIS_URL ${REDIS_URL}

COPY --from=builder /usr/src/app/nodes /home/nodes

CMD ["./nodes"]
