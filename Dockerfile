FROM golang:1.16.4-alpine3.13 AS builder

RUN apk add --update --no-cache git ca-certificates && update-ca-certificates

WORKDIR /go/src/app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o adak -ldflags="-s -w" ./cmd/main.go

# --------------------------------------------

FROM scratch

COPY --from=builder /go/src/app/adak /usr/bin/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/usr/bin/adak"]