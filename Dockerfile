FROM golang:1.23.3-alpine AS builder

ARG GOOS=linux
ARG GOARCH=amd64
ARG CGO_ENABLED=1
ARG GOFLAGS="-ldflags=-extldflags=-static -a -v"

WORKDIR /code

RUN apk add \
    gcc \
    musl-dev \
    libpcap-dev

COPY go.mod go.sum ./
RUN go mod download -x

COPY main.go .
COPY internal internal

RUN go build

FROM scratch AS final

COPY --from=builder /code/listener .
ENTRYPOINT ["/listener"]

