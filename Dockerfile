# Builder image
FROM golang:1.20 AS builder
ARG BUILD_VERSION=docker-dev

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/
RUN CGO_ENABLED=0 go build \
      -o /server \
      ./cmd/server/*.go

# We want to access some basic shell tools for debugging, but we want to be
# as tiny as possible...
FROM alpine:3.17.0
RUN apk add strace
COPY --from=builder /server /usr/local/bin/server

ENTRYPOINT ["/usr/local/bin/server"]
