FROM golang:1.23.3-alpine AS builder

ARG GOOS=linux
ARG GOARCH=amd64
ARG CGO_ENABLED=1

WORKDIR /code

RUN apk add \
    gcc \
    musl-dev \
    libpcap-dev

COPY go.mod go.sum ./
RUN go mod download -x

COPY main.go .
COPY cmd cmd
COPY internal internal
RUN go build
RUN go build -o harkener -tags netgo,osusergo -ldflags '-extldflags "-static" -w -s' .

FROM scratch AS final

COPY --from=builder /code/harkener .
ENTRYPOINT ["/harkener"]

