[![Build Status](https://app.travis-ci.com/kutzilla/docker-hetzner-ddns.svg?branch=master)](https://app.travis-ci.com/kutzilla/docker-hetzner-ddns) [![Docker Pulls](https://img.shields.io/docker/pulls/kutzilla/hetzner-ddns.svg)](https://hub.docker.com/r/kutzilla/hetzner-ddns)

# Docker Hetzner DDNS

This Docker image will allow you to use the [Hetzner DNS Service](https://www.hetzner.com/dns-console) as a Dynamic DNS Provider ([DDNS](https://en.wikipedia.org/wiki/Dynamic_DNS)).


## Usage

Quick Setup:

```shell
docker run \
-e ZONE_NAME=example.com \ 
-e HETZNER_API_TOKEN=my-secret-api-token \
-e ZONE_RECORD_TYPE=A \
kutzilla/hetzner-ddns
```


If you prefer command-line arguments, you can use this command: 

```shell
docker run kutzilla/hetzner-ddns example.com my-secret-api-token A
```

## Parameters


* `-e ZONE_NAME` - The DNS zone that DDNS updates should be applied to. **Required**
* `-e API_TOKEN` - Your Hetzner  API token. **Required**
* `-e RECORD_TYPE` - The record type of your zone. If your zone uses an IPv4 address use `A`. Use `AAAA` if it uses an IPv6 address. **required**
* `--restart=always` - ensure the container restarts automatically after host reboot.



## Build

Build the latest version of the Docker image with the following command:

```
docker build kutzilla/hetzner-ddns .
```

