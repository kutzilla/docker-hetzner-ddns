[![Build Status](https://app.travis-ci.com/kutzilla/docker-hetzner-ddns.svg?branch=master)](https://app.travis-ci.com/kutzilla/docker-hetzner-ddns)

# Docker Hetzner DDNS

This Docker image will allow you to use the [Hetzner DNS Service](https://www.hetzner.com/dns-console) as a Dynamic DNS Provider ([DDNS](https://en.wikipedia.org/wiki/Dynamic_DNS)).


## Usage


### Docker

You can run the Docker image with the following command:

```
docker run kutzilla/hetzner-ddns {ZONE_NAME} {HETZNER_API_TOKEN} {ZONE_RECORD_TYPE}
```

If you prefer environment variables, you can use this command: 

```
docker run -e ZONE_NAME=example.com -e HETZNER_API_TOKEN=my-secret-api-token -e ZONE_RECORD_TYPE=A kutzilla/hetzner-ddns
```

### Go

You also can can run the Go implementation with the following command:

```
go run main.go {ZONE_NAME} {HETZNER_API_TOKEN} {ZONE_RECORD_TYPE}
```

If you prefer environment variables, you can use this command: 

```
set ZONE_NAME=example.com
set HETZNER_API_TOKEN=my-secret-api-token
set ZONE_RECORD_TYPE=A // or AAAA

go run main.go
```


# Build

Build the latest version of the Docker image with the following command:

```
docker build kutzilla/hetzner-ddns .
```

