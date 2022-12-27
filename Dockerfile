FROM golang:1.19-alpine AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY pkg ./pkg
COPY cmd ./cmd
RUN go build -o hetzner-ddns ./cmd/hetzner-ddns

FROM scratch
WORKDIR /
COPY --from=build /app/hetzner-ddns /hetzner-ddns
ENTRYPOINT ["/hetzner-ddns"]
