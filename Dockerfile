FROM golang:1.16.0-alpine3.13 AS builder

WORKDIR /go/src/app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o adak -ldflags="-s -w" ./cmd/main.go

# -------------------------

FROM scratch

COPY --from=builder /go/src/app/adak /usr/bin/

ENTRYPOINT ["/usr/bin/adak"]