FROM golang:1.15-alpine AS builder

COPY . /go/src/app

WORKDIR /go/src/app

RUN go get -d -v ./...

RUN go build -o adak -ldflags="-s -w" ./cmd/main.go

# -------------------------

FROM alpine:3.13.2

COPY --from=builder /go/src/app/adak /usr/bin/

ENTRYPOINT ["/usr/bin/adak"]