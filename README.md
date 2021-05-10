# Adak

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/GGP1/adak/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/GGP1/adak)](https://goreportcard.com/report/github.com/GGP1/adak)

Adak is an e-commerce server developed for educational purposes.

This project has two versions: **monolithic** and **microservices**.

For the microservices version please change to the [microservices](https://github.com/GGP1/adak/tree/microservices) branch.

## Features

- Authentication
    - Password encryption using [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
    - Email verification
    - Distinction between users and admins
    - Delay after a configurable amount of failures (for mitigating brute force attacks)
    - Basic authentication and OAuth2 (Google)
- Distributed caching, sessions and rate limiting
- Rate limiter to protect against API abuse and DDoS attacks
- Hexagonal architecture
- Context cancelling and graceful shutdown
- OpenAPI Specification 3.0.0 with Swagger
- Requests and services logging
- Connection encrypted with TLS (self-signed certificate as it wont be used in production)
- Input sanitization and validation
- Encrypted cookies using ChaCha20-Poly1305 and hex encoding/decoding
- Non invasive user tracking service
- GZIP responses compression
- Stripe integration for payments

## Installation

### Docker

> Optional: mount your configuration files and certificates in the [docker-compose.yml](/docker-compose.yml) file, or remove them to use the default configuration.

```
docker-compose up
```

### Manually

> Requires postgres, memcached, redis, Go and Git to be installed

Clone the repository: 
```
git clone https://www.github.com/GGP1/adak.git
```

**Configuration**: set an environment variable called `ADAK_CONFIG` pointing to your [configuration file](/config_example.yml).

Run the server: 

```bash
go run cmd/main.go
# or
make run
```

## Documentation

### API

The API's documentation can be found on [SwaggerHub](https://app.swaggerhub.com/apis/GGP1/ADAK_OAS3/1.0.0).

### Database migrations

The [golang-migrate/migrate](https://github.com/golang-migrate/migrate) CLI is required to execute the commands.

Up:
```
migrate -path pkg/postgres/migrations -database <postgres_url> up
```

Down:
```
migrate -path pkg/postgres/migrations -database <postgres_url> down
```

### Metrics

Adak collects information using [prometheus](https://prometheus.io/) and runs a [grafana](https://grafana.com/) container for visualizing it.

The data is gathered from multiple sources:

- **Server**: requests and database hits
- **Node**: hardware and OS metrics provided by [node exporter](https://github.com/prometheus/node_exporter)

More data can be extracted from other services like Postgres, Redis and Memcached, although they are not implemented, it would imply configuration changes only.

#### Visualizing data

There are a few ways of see the actual data:

- The server raw statistics are showed in the */metrics* path (it also shows Go's metrics but not the information from external sources).
- With prometheus by creating graphs at `http://localhost:9090` 
- Grafana

To display prometheus' information with grafana follow these steps:

1. Go to http://localhost:3000 and login as "admin" user, the password is "admin".
2. Add prometheus as a data source, set `HTTTP/Access` to `Browser` and the url to `http://localhost:9090`, then click `Save & Test`.
3. On the top of the page and next to the *Settings* tab there other called *Dashboards*, click it and import `Prometheus 2.0 Stats`.
4. Finally, go to the dashboards home and select the dashboard imported above.

> You can create your own or use official/community built [dashboards](https://grafana.com/grafana/dashboards). This is more related to Grafana than Adak so it's left to the user.

### Amounts

Amounts are represented by 64-bit integers to be provided in a currency's smallest unit (100 = 1 USD).

In the case of weights, 1000 = 1kg.
