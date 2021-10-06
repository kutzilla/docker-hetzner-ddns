# Docker Hetzner DDNS

This Docker image will allow you to use the [Hetzner DNS Service](https://www.hetzner.com/dns-console) as a Dynamic DNS Provider ([DDNS](https://en.wikipedia.org/wiki/Dynamic_DNS)).

## Usage

### Go

```
go run main.go {ZONE_NAME} {HETZNER_API_TOKEN} {ZONE_RECORD_TYPE}
```

### Docker

```
docker build -t matthias-kutz.com/hetzner-ddns .
docker run matthias-kutz.com/hetzner-ddns {ZONE_NAME} {HETZNER_API_TOKEN} {ZONE_RECORD_TYPE}
```


