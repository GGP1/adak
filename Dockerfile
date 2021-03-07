FROM golang:1.16.0-alpine3.13 AS builder

WORKDIR /go/src/app

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o adak -ldflags="-s -w" ./cmd/main.go

# -------------------------

FROM alpine:3.13.2

COPY --from=builder /go/src/app/adak /usr/bin/

ENTRYPOINT ["/usr/bin/adak"]