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

### Amounts

Amounts are represented by 64-bit integers to be provided in a currency's smallest unit (100 = 1 USD).

In the case of weights, 1000 = 1kg.
