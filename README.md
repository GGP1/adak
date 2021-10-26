# Adak

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/GGP1/adak/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/GGP1/adak)](https://goreportcard.com/report/github.com/GGP1/adak)

Adak is an e-commerce server developed for educational purposes.

This project has two versions: **monolithic** and **microservices**.

For the microservices version please change to the [microservices](https://github.com/GGP1/adak/tree/microservices) branch.

## Features

- Cookie-based sessions encrypted using ChaCha20-Poly1305
- Basic authentication and OAuth2 (Google)
- Password encryption using [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- Email and admins verification
- OpenAPI Specification 3.0.0 with Swagger
- Stripe integration for payments
- Pagination, caching, rate limiting, GZIP responses compression, input sanitization and validation, context cancelling

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

### Pagination

Adak accepts bot `limit` and `cursor` url parameters for paginating results. 

The `next_cursor` field is always included in the response if there are more values to be fetched. The next cursor is builded up from the last record's *created_at timestamp* and *uuid*, they are formatted and concatenated, the resulting string is encoded with base64 and sent to the client.

Example response:

```json
{
"next_cursor": "MjAyMS0wNy0wOVQyMzo0ODozOC42NDg0MDZaLHhLM0lBY01XenFTdmY5aEdyUVpsOENOc1haeUV2NlBJ",
    "users": [
        {
            "id": "xK3IAcMWzqSvf9hGrQZl8CNsXZyEv6PI",
            "created_at": "2021-07-09T23:48:38.648406Z",
        }
    ]
}
```

Example request:
```
localhost:4000/users?limit=5&cursor=MjAyMS0wNy0wOVQyMzo0ODozOC42NDg0MDZaLHhLM0lBY01XenFTdmY5aEdyUVpsOENOc1haeUV2NlBJ
```

### Amounts

Amounts are represented by 64-bit integers to be provided in a currency's smallest unit (100 = 1 USD).

In the case of weights, 1000 = 1kg.

### Monitoring

Adak collects information using [prometheus](https://prometheus.io/) and runs a [grafana](https://grafana.com/) container for visualizing it.

The data is gathered from multiple sources:

- **Server**: requests and database hits
- **Node**: hardware and OS metrics provided by [node exporter](https://github.com/prometheus/node_exporter)

More data can be extracted from other services like Postgres, Redis and Memcached, although they are not implemented, it would imply configuration changes only.

### Visualizing metrics

#### Raw

The server's raw statistics are shown in the "*/metrics*" path (it also contains Go metrics but not information from external sources).

#### Prometheus

Prometheus allows not only querying the data but also displaying it in graphs, this can be done at `http://localhost:9090` when running with docker compose.

![Prometheus](https://user-images.githubusercontent.com/51374959/118064036-a459f500-b370-11eb-999b-6e539c5b4b9f.png)

#### Grafana

To display prometheus' metrics with grafana follow these steps:

1. Go to http://localhost:3000 and log in as "admin", the password is "admin".
2. Add prometheus as a data source. Set `HTTTP/Access` to `Browser` and the url to `http://localhost:9090`, then click `Save & Test`.
3. On the top of the page and next to the *Settings* tab there other called *Dashboards*, click it and import `Prometheus 2.0 Stats`.
4. Finally, go to the dashboards home and select the dashboard imported above.

> You can create your own or use official/community built [dashboards](https://grafana.com/grafana/dashboards).

![Grafana](https://user-images.githubusercontent.com/51374959/118064057-ade35d00-b370-11eb-9fc2-4fa2dc859c8b.png)

### Load testing

Making 10000 requests with 100 concurrent workers to the `/users` endpoint, each response containing 20 non-cached user objects (2.94KB per response) resulted in:

```
$ hey -n 10000 -c 100 -m GET http://:4000/users

Summary:
  Total:        11.4822 secs
  Slowest:      0.4250 secs
  Fastest:      0.0021 secs
  Average:      0.1110 secs
  Requests/sec: 870.9148

Response time histogram:
  0.002 [1]     |
  0.044 [1708]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.087 [2219]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.129 [2383]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.171 [1834]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.214 [1141]  |■■■■■■■■■■■■■■■■■■■
  0.256 [491]   |■■■■■■■■
  0.298 [169]   |■■■
  0.340 [38]    |■
  0.383 [13]    |
  0.425 [3]     |

Latency distribution:
  10% in 0.0281 secs
  25% in 0.0608 secs
  50% in 0.1049 secs
  75% in 0.1542 secs
  90% in 0.1982 secs
  95% in 0.2284 secs
  99% in 0.2809 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0021 secs, 0.4250 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0041 secs
  resp wait:    0.1107 secs, 0.0019 secs, 0.4247 secs
  resp read:    0.0002 secs, 0.0000 secs, 0.0064 secs

Status code distribution:
  [200] 10000 responses
```
