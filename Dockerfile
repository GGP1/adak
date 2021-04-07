FROM golang:1.16.2-alpine3.13 AS builder

WORKDIR /microservices

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o adak -ldflags="-s -w" ./cmd/main.go

# --------------------------

FROM scratch

COPY --from=builder /microservices/adak /usr/bin/adak

ENTRYPOINT ["/usr/bin/adak"]