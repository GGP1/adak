FROM golang:1.14-alpine AS builder

########
# Prep
########

# add the source
COPY . /go/src/palo
WORKDIR /go/src/palo/

########
# Build Go Wrapper
########

# Install go dependencies
RUN go get -d -v ./...

#build the go app
RUN GOOS=linux GOARCH=amd64 go build -o ./palo ./cmd/main.go

########
# Package into runtime image
########
FROM alpine

# copy the executable from the builder image
COPY --from=builder /go/src/palo .

ENTRYPOINT ["/palo"]

EXPOSE 8080