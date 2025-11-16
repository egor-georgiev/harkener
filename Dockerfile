ARG OS=linux
ARG ARCH=amd64

FROM --platform=${OS}/${ARCH} golang:1.25.4-alpine AS builder
ARG OS
ARG ARCH
ARG GOOS=${OS}
ARG GOARCH=${ARCH}
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
RUN go build -o harkener-${OS}-${ARCH} -tags netgo,osusergo -ldflags '-extldflags "-static" -w -s' .

FROM scratch AS final
ARG OS
ARG ARCH

COPY --from=builder /code/harkener-${OS}-${ARCH} .
ENTRYPOINT ["/harkener"]

