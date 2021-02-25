FROM golang:1.16.0-alpine3.13 AS builder

COPY . /microservices
WORKDIR /microservices

RUN CGO_ENABLED=0 go build -o adak -ldflags="-s -w" ./cmd/main.go

# --------------------------

FROM scratch

COPY --from=builder /microservices/adak /bin/adak

# ENTRYPOINT ["/bin/adak"]