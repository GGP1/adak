# Adak

[![PkgGoDev](https://pkg.go.dev/badge/github.com/GGP1/adak)](https://pkg.go.dev/github.com/GGP1/adak)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/GGP1/adak/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/GGP1/adak)](https://goreportcard.com/report/github.com/GGP1/adak)

Adak is an e-commerce server developed for educational purposes.

This project has two versions: **monolithic** and **microservices**.

For the microservices version please change to the [microservices](https://github.com/GGP1/adak/tree/microservices) branch.

## Features

- Authentication
    - Password encryption using [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
    - Email verification
    - Distinction between users and admins
    - Incremental delay when login from the same ip fails multiple times (for mitigating brute force attacks)
- Rate limiter to protect against API abuse and DDoS attacks (DDoSing is stil possible but it will require a high number of computers faking handshakes in parallel)
- Context cancelling
- API documentation using Swagger
- Requests and services logging
- Graceful shutdown
- Connection encrypted with TLS (self-signed certificate as it wont be used in production)
- Input sanitization and unicode normalization to NFC
- Encrypted cookies using ChaCha20-Poly1305 and hex encoding/decoding
- Non invasive user tracking service
- GZIP responses compression (if accepted by the client)
- Stripe integration
- Email and password changing
- Google OAuth2

## Installation

#### Docker

*Build images and run containers*: `docker-compose up`

The server will be listening on port **4000** and the documentation on port **8080**.

#### Manually

> Requires postgres server to be installed

*Clone the repository*: `git clone https://www.github.com/GGP1/adak.git`

*Configuration*: set an environment variable called "ADAK_CONFIG" pointing to your [.env file](/.env_example).

*Run the server*: `go run cmd/main.go`

## Documentation

Please head over to the [API docs](/api/swagger/README.md) for further details.

### Endpoints

|  Service  |  Method   |  Auth   |              URL                  |            Description            |
|-----------|-----------|---------|-----------------------------------|-----------------------------------|
|  Auth     | POST      | None    | /login                            | Session login                     |
|           | GET       | Login   | /logout                           | Session logout                    |
|           | GET       | None    | /login/google                     | Redirect to the google oauth login|
|           | GET       | None    | /login/oauth2/google              | Session login with Google         |
|-----------|-----------|---------|-----------------------------------|-----------------------------------|
|  Cart     | GET       | Login   | /cart                             | Get cart details                  |
|           | POST      | Login   | /cart/add/{quantity}              | Add a product to the cart         |
|           | GET       | Login   | /cart/brand/{brand}               | Filter cart products by brand     |
|           | GET       | Login   | /cart/category/{category}         | Filter cart products by category  |
|           | GET       | Login   | /cart/discount/{min}/{max}        | Filter cart products by discount  |
|           | GET       | Login   | /cart/checkout                    | Fetch cart total                  |
|           | GET       | Login   | /cart/products                    | Fetch cart products               |
|           | DELETE    | Login   | /cart/remove/{id}/{quantity}      | Delete a product from the cart    |
|           | GET       | Login   | /cart/reset                       | Set cart fields to default        |
|           | GET       | Login   | /cart/size                        | Get cart number of products       |
|           | GET       | Login   | /cart/taxes/{min}/{max}           | Filter cart products by taxes     |
|           | GET       | Login   | /cart/total/{min}/{max}           | Filter cart products by total     |
|           | GET       | Login   | /cart/type/{type}                 | Filter cart products by type      |
|           | GET       | Login   | /cart/weight/{min}/{max}          | Filter cart products by weight    |
|-----------|-----------|---------|-----------------------------------|-----------------------------------|
|  Home     | GET       | None    | /                                 | Adak welcome page                 |
|-----------|-----------|---------|-----------------------------------|-----------------------------------|
|  Order    | GET       | Admin   | /orders                           | Get all orders                    |
|           | DELETE    | Admin   | /order/{id}                       | Delete an order                   |
|           | GET       | Admin   | /order/{id}                       | Get order by id                   |
|           | GET       | Login   | /order/user/{id}                  | Get orders from a user            |
|           | POST      | Login   | /order/new                        | Create new order                  |
|-----------|-----------|---------|-----------------------------------|-----------------------------------|
|  Product  | GET       | None    | /products                         | Get all products                  |
|           | GET       | None    | /products/{id}                    | Get product by id                 |
|           | PUT       | Admin   | /products/{id}                    | Update product                    |
|           | DELETE    | Admin   | /products/{id}                    | Delete a product                  |
|           | POST      | Admin   | /products/create                  | Create a product                  |
|           | GET       | None    | /products/search/{query}          | Products search                   |
|-----------|-----------|---------|-----------------------------------|-----------------------------------|
|  Review   | GET       | None    | /reviews                          | Get all reviews                   |
|           | GET       | None    | /reviews/{id}                     | Get review by id                  |
|           | DELETE    | Admin   | /reviews/{id}                     | Delete a review                   |
|           | PUT       | Admin   | /reviews/{id}                     | Update review                     |
|           | POST      | Login   | /reviews/create                   | Create a review                   |
|-----------|-----------|---------|-----------------------------------|-----------------------------------|
|  Shop     | GET       | None    | /shops                            | Get all shops                     |
|           | GET       | None    | /shops/{id}                       | Get shop by id                    |
|           | DELETE    | Admin   | /shops/{id}                       | Delete a shop                     |
|           | PUT       | Admin   | /shops/{id}                       | Update shop                       |
|           | POST      | Admin   | /shops/create                     | Create a shop                     |
|           | GET       | None    | /shops/search/{query}             | Shops search                      |
|-----------|-----------|---------|-----------------------------------|-----------------------------------|
|  Stripe   | GET       | Admin   | /stripe/balance                   | Get stripe balance                |
|           | GET       | Admin   | /stripe/event/{event}             | Get stripe event                  |
|           | GET       | Admin   | /stripe/transactions/{txID}       | Get transaction details           |
|           | GET       | Admin   | /stripe/events                    | Get all stripe events             |
|           | GET       | Admin   | /stripe/transactions              | Get all stripe transactions       |
|-----------|-----------|---------|-----------------------------------|-----------------------------------|
|  Tracker  | GET       | Admin   | /tracker                          | Get tracker Hits                  |
|           | GET       | Admin   | /tracker/{id}                     | Get tracket hit by id             |
|           | GET       | Admin   | /tracker/search/{query}           | Hits search                       |
|           | GET       | Admin   | /tracker/{field}/{value}          | Hits search by field              |
|-----------|-----------|---------|-----------------------------------|-----------------------------------|
|  Users    | GET       | None    | /users                            | Get all users                     |
|           | GET       | None    | /users/{id}                       | Get user by id                    |
|           | PUT       | Login   | /users/{id}                       | Update user                       |
|           | DELETE    | Admin   | /users/{id}                       | Delete user                       |
|           | POST      | None    | /users/create                     | Create new user                   |
|           | GET       | None    | /users/search/{search}            | Users search                      |
|-----------|-----------|---------|-----------------------------------|-----------------------------------|
| Account   | POST      | Login   | /settings/email                   | Send email change confirmation    |
|           | POST      | Login   | /settings/password                | Change account password           |
|           | GET       | None    | /verification/{email}/{token}     | Send email validation             |
|           | GET       | None    | /verification/{token}/{email}/{id}| Change account email              |

### Amounts

Amounts are represented by 64-bit integers to be provided in a currency’s smallest unit (100 = 1 USD).

In the case of weights, 1000 = 1kg.
