# Replace the release candidate once the 1.16 image is published
FROM golang:1.16rc1-alpine3.13 AS builder

COPY . /adak

WORKDIR /adak

RUN CGO_ENABLED=0 go build -o api -ldflags="-s -w" .

# -------------------------

FROM scratch

COPY --from=builder /adak/api /GGP1/adak/

EXPOSE 8080

ENTRYPOINT ["/GGP1/adak/api"]