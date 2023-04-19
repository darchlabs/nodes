## 1st stage: build golang ethereum runner
FROM golang:alpine as builder

WORKDIR /usr/src/app

COPY . .

RUN CGO_ENABLED=0 go build -o nodes cmd/nodes/main.go

## 2nd stage: prepare container to run node
FROM golang:alpine as runner

WORKDIR /home/nodes

COPY --from=builder /usr/src/app/nodes /home/nodes

CMD ["./nodes"]
