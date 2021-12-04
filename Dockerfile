FROM golang:1.17-alpine3.14

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY /pkg /app/pkg
COPY  /cmd /app/cmd

RUN go mod download
RUN go build -o hetzner-ddns ./cmd/hetzner-ddns

ENTRYPOINT ["/app/hetzner-ddns"]