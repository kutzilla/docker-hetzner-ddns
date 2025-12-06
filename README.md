[![Build Status](https://github.com/kutzilla/docker-hetzner-ddns/actions/workflows/build.yml/badge.svg)](https://github.com/kutzilla/docker-hetzner-ddns/actions/workflows/go.yml) [![Docker Pulls](https://img.shields.io/docker/pulls/kutzilla/hetzner-ddns.svg)](https://hub.docker.com/r/kutzilla/hetzner-ddns)

# Docker Hetzner DDNS

This Docker image will allow you to use the [Hetzner DNS Service](https://www.hetzner.com/dns-console) as a Dynamic DNS Provider ([DDNS](https://en.wikipedia.org/wiki/Dynamic_DNS)).

## How does it work?

The Go script inside this Docker Image periodically checks the DNS record with the Hetzner DNS API. It also checks the current public IP of the network, the container is running on. If the DNS record does not match the current public IP, it will update the record. Therefore your DNS record updates dynamically to the public IP.



## Usage

Quick Setup:

```shell
docker run \
-e ZONE_NAME=example.com \ 
-e API_TOKEN=my-secret-api-token \
-e RECORD_TYPE=A \
-e RECORD_NAME_VPN=vpn \
-e RECORD_NAME_VPN_TTL=600 \
-e RECORD_NAME_WEBSITE=www \
-e RECORD_NAME_WEBSITE_TTL=1337 \
kutzilla/hetzner-ddns
```


If you prefer command-line arguments, you can use this command: 

```shell
docker run kutzilla/hetzner-ddns example.com my-secret-api-token A
```

## Parameters


* `-e ZONE_NAME` - The DNS zone that DDNS updates should be applied to. **Required**
* `-e API_TOKEN` - Your Hetzner API token. **Required**
* `-e RECORD_TYPE` - The record type of your zone. If your zone uses an IPv4 address use `A`. Use `AAAA` if it uses an IPv6 address. **Required**
* `--restart=always` - ensure the container restarts automatically after host reboot.

## Optional Parameters

* `-e RECORD_NAME` - The name of the DNS-record that DDNS updates should be applied to. This could be `sub` if you like to update the subdomain `sub.example.com` of `example.com`. The default value is `@` If you want to update multiple Records you can use a pattern. The pattern which can be used is `RECORD_NAME_<NAME>` e.g. `RECORD_NAME_EXAMPLE`, `RECORD_NAME` will be ignored the multi domain approach is used.
* `-e CRON_EXPRESSION` - The cron expression of the DDNS update interval. The default is every 5 minutes - `*/5 * * * *`.
* `-e TTL` - The TTL (Time To Live) value specifies how long the record is valid before the nameservers are prompted to reload the zone file. The default is `86400`.
* `-e RECORD_NAME_<NAME>_TTL` - In case you have multiple records, you can specify the TTL per record referenced by `<NAME>` for example `RECORD_NAME_EXAMPLE_TTL`, if no TTL is specified the TTL will be `86400`. The `TTL` argument will be ignored if this parameter is used.


## Build

Build the latest version of the Docker image with the following command:

```
docker build -t kutzilla/hetzner-ddns .
```

