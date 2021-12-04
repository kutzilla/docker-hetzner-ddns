[![Build Status](https://app.travis-ci.com/kutzilla/docker-hetzner-ddns.svg?branch=master)](https://app.travis-ci.com/kutzilla/docker-hetzner-ddns) [![Docker Pulls](https://img.shields.io/docker/pulls/kutzilla/hetzner-ddns.svg)](https://hub.docker.com/r/kutzilla/hetzner-ddns)

# Docker Hetzner DDNS

This Docker image will allow you to use the [Hetzner DNS Service](https://www.hetzner.com/dns-console) as a Dynamic DNS Provider ([DDNS](https://en.wikipedia.org/wiki/Dynamic_DNS)).

## How does it work?

The Go script inside this Docker Image periodically checks the DNS record with the Hetzner DNS API. It also checks the current public IP of the network, the container is running on. If the DNS record does not match the current public IP, it will update the record. Therefore your DNS record updates dynamically to the public IP.



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
* `-e RECORD_TYPE` - The record type of your zone. If your zone uses an IPv4 address use `A`. Use `AAAA` if it uses an IPv6 address. **Required**
* `--restart=always` - ensure the container restarts automatically after host reboot.

## Optional Parameters

* `-e RECORD_NAME` - The name of the DNS-record that DDNS updates should be applied to. This could be `sub` if you like to update the subdomain `sub.example.com` of `example.com`. The default value is `@`.
* `-e CRON_EXPRESSION` - The cron expression of the DDNS update interval. The default is every 5 minutes - `*/5 * * * *`.

## Build

Build the latest version of the Docker image with the following command:

```
docker build -t kutzilla/hetzner-ddns .
```

