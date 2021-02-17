## Documentation

The API documentation is built by Swagger, display it by doing:

### Docker-compose

1. Go to the repository root
2. `docker-compose up`
3. Open *localhost* on your browser:
- *:4000* for the server
- *:8080* for the documentation

### Docker (Remote) 

1. `docker run --rm -p 8080:8080 gastonpalomeque/adak-api`
2. Open *localhost:8080* on your browser

### Docker (Local)

1. `docker build -t adak-api .`
2. `docker run --rm -p 8080:8080 adak-api`
3. Open *localhost:8080* on your browser

### Go

1. `go run api/swagger/main.go --port 4000`
2. Open *localhost:4000* on your browser

> Default port is 8080.

### Manually

1. Open the *index.html* file with your browser.