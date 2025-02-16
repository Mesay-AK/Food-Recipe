FROM golang:1.22.5 AS builder

WORKDIR /backend

COPY go.mod go.sum ./

RUN go mod download

COPY . . 

WORKDIR /backend/cmd

RUN go build -o /backend/main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /backend/main .

EXPOSE 8081

CMD ["./main"]
